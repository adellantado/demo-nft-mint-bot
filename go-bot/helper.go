package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/progress"
)

func GetChatId(update *models.Update) int64 {
	return update.Message.Chat.ID
}

func GetTextMsgParams(update *models.Update, text string) *bot.SendMessageParams {
	return &bot.SendMessageParams{
		ChatID: GetChatId(update),
		Text:   text,
	}
}

func HasPhoto(update *models.Update) ([]models.PhotoSize, bool) {
	photo := update.Message.Photo
	if len(photo) > 0 {
		return photo, true
	}
	return photo, false
}

func InitProgressBar(ctx context.Context, b *bot.Bot, update *models.Update) (*progress.Progress, *time.Ticker, chan bool) {
	opts := []progress.Option{
		progress.WithRenderTextFunc(func(value float64) string {
			return bot.EscapeMarkdown(fmt.Sprintf("Minting progress: %.2f%%", value))
		}),
	}
	p := progress.New(b, opts...)
	p.Show(ctx, b, GetChatId(update))

	ticker := time.NewTicker(1000 * time.Millisecond)
	ticker_count := 0.00
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				ticker_count++
				p.SetValue(ctx, b, ticker_count*5)
				if ticker_count >= 100.00 {
					ticker.Stop()
					done <- true
				}
			}
		}
	}()

	return p, ticker, done
}

func DeleteProgressBar(ctx context.Context, b *bot.Bot, p *progress.Progress, ticker *time.Ticker, done chan bool) {
	ticker.Stop()
	done <- true
	p.Delete(ctx, b)
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
