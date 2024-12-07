package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/joho/godotenv"
)

var storage *UserStorage
var convoManager *ConversationManager

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	storage = NewUserStorage()
	convoManager = NewConversationManager()
	convoManager.AddConvoHandlers(map[string][]func(context.Context, *bot.Bot, *models.Update) string{
		"mintConvo": {
			imageMintConvoHandler,
			titleMintConvoHandler,
			descriptionMintConvoHandler,
		},
	})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/mint", bot.MatchTypeExact, mintHandler)

	fmt.Println("The bot is running! Press Ctrl+C to terminate!")
	b.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, GetTextMsgParams(update, "Hello! Please run /mint command"))
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if convoManager.Handle(ctx, b, update) {
		return
	}

	b.SendMessage(ctx, GetTextMsgParams(update, "Sorry, I don't understand.. but I can run /mint command for you"))
}

func mintHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, GetTextMsgParams(update, "Please upload an image for new NFT"))
	convoManager.InitConvo(GetChatId(update), "mintConvo")
}

func imageMintConvoHandler(ctx context.Context, b *bot.Bot, update *models.Update) string {
	photo, exists := HasPhoto(update)
	if exists {
		ps := photo[len(photo)-1]
		storage.Store(update, "imageId", ps.FileID)
		b.SendMessage(ctx, GetTextMsgParams(update, "Perfect! The image is loaded.\nType a title for the image below:"))
		return "1"
	}
	b.SendMessage(ctx, GetTextMsgParams(update, "Error! Please send image instead"))
	return "0"
}

func titleMintConvoHandler(ctx context.Context, b *bot.Bot, update *models.Update) string {
	storage.Store(update, "title", update.Message.Text)
	b.SendMessage(ctx, GetTextMsgParams(update, "Before we proceed with minting, please, provide some short NFT description below:"))
	return "2"
}

func descriptionMintConvoHandler(ctx context.Context, b *bot.Bot, update *models.Update) string {
	storage.Store(update, "description", update.Message.Text)

	fileId := storage.Get(update, "imageId")
	title := storage.Get(update, "title")
	description := storage.Get(update, "description")
	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: fileId})
	if err != nil {
		fmt.Println(err)
		b.SendMessage(ctx, GetTextMsgParams(update, "The error with image is occured! Please try again /mint"))
		return END
	}
	fileUrl := b.FileDownloadLink(file)
	fileExt := filepath.Ext(fileUrl)
	fileName := fileId + fileExt

	err = DownloadFile("images/"+fileName, fileUrl)
	if err != nil {
		fmt.Println(err)
		b.SendMessage(ctx, GetTextMsgParams(update, "The error with image downloading is occured! Please try again /mint"))
		return END
	}

	b.SendMessage(ctx, GetTextMsgParams(update, "Initializing minting, wait a sec..."))

	mint_api_url := os.Getenv("MINTING_ENDPOINT")
	json_data := map[string]string{
		"title":       title,
		"description": description,
		"image":       fileName,
	}

	bytes_data, _ := json.Marshal(json_data)
	req, err := http.NewRequest("POST", mint_api_url, bytes.NewBuffer(bytes_data))
	if err != nil {
		b.SendMessage(ctx, GetTextMsgParams(update, "Error with the data for minting endpoint! Please try again /mint"))
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	p, ticker, done := InitProgressBar(ctx, b, update)

	resp, err := client.Do(req)

	DeleteProgressBar(ctx, b, p, ticker, done)
	if err != nil {
		b.SendMessage(ctx, GetTextMsgParams(update, "Error with minting endpoint! Please try again /mint"))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	res_data := map[string]any{
		"title":       "",
		"description": "",
		"image":       "",
		"mint":        "",
	}
	err = json.Unmarshal(body, &res_data)
	if err != nil {
		fmt.Println(err)
		b.SendMessage(ctx, GetTextMsgParams(update, "Error with response decoding! Please try again /mint"))
	}
	b.SendMessage(ctx, GetTextMsgParams(update, "Mint link: "+res_data["mint"].(string)))

	return END
}
