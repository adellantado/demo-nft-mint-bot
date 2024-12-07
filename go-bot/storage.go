package main

import "github.com/go-telegram/bot/models"

type UserStorage struct {
	userdata map[int64]map[string]string
}

func (us *UserStorage) Store(update *models.Update, key string, value string) {
	chatId := update.Message.Chat.ID
	if _, ok := us.userdata[chatId]; !ok {
		us.userdata[chatId] = make(map[string]string)
	}
	us.userdata[chatId][key] = value
}

func (us *UserStorage) Get(update *models.Update, key string) string {
	return us.userdata[update.Message.Chat.ID][key]
}

func (us *UserStorage) StoreByChatId(chatId int64, key string, value string) {
	if _, ok := us.userdata[chatId]; !ok {
		us.userdata[chatId] = make(map[string]string)
	}
	us.userdata[chatId][key] = value
}

func (us *UserStorage) GetByChatId(chatId int64, key string) string {
	return us.userdata[chatId][key]
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		userdata: make(map[int64]map[string]string),
	}
}
