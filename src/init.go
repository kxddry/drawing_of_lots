package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"strconv"
)

const (
	MaxUsers         = 28
	NumberOfGroups   = 3
	MaxUsersPerGroup = 10 // will have to check back on that
)

var (
	envFile, errOpen     = godotenv.Read("pkg/.env")          // the .env file is located in the pkg folder
	token                = envFile["TOKEN"]                   // the telegram bot API token
	owner                = envFile["OWNER"]                   // the ID of the one starting the polls
	int64owner, errParse = strconv.ParseInt(owner, 10, 64)    // 64 bytes is required for userIDs
	peers                = make([]int64, 0, MaxUsers)         // users
	usersHashmap         = make(map[int64][]string, MaxUsers) // {...chatID: [username, firstName]...}
	groups               = genGroups()
	punishUser           = false
)

// keyboards
var (
	noPollKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Список"),
			tgbotapi.NewKeyboardButton("Помощь"),
		),
	)
	ownerKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Список"),
			tgbotapi.NewKeyboardButton("Помощь"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Poll"),
			tgbotapi.NewKeyboardButton("Shutdown"),
		),
	)
	pollOwnerKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Poll"),
			tgbotapi.NewKeyboardButton("Shutdown"),
		),
	)
)

// inline keyboards

var (
	noPollInline = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Регистрация", "register"),
			tgbotapi.NewInlineKeyboardButtonData("Выход", "quit"),
		),
	)
	registerInline = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Регистрация", "register"),
		),
	)
	quitInline = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Выход", "quit"),
		),
	)
	sendInline = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить", "send"),
		),
	)
)

func init() {
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	if token == "" {
		log.Fatal("check your API token")
	}
	if owner == "" {
		log.Fatal("No owner provided")
	}
	if errParse != nil {
		log.Fatal(errParse)
	}
	// if any issues during reading the env file arise, stop the program
}
