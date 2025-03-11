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
	BotLanguage          = envFile["language"]
	envFile, errOpen     = godotenv.Read("pkg/.env")          // the .env file is located in the pkg folder
	token                = envFile["TOKEN"]                   // the telegram bot API token
	owner                = envFile["OWNER"]                   // the ID of the one starting the polls
	int64owner, errParse = strconv.ParseInt(owner, 10, 64)    // 64 bytes is required for userIDs
	peers                = make([]int64, 0, MaxUsers)         // users
	usersHashmap         = make(map[int64][]string, MaxUsers) // {...chatID: [username, firstName]...}
	groups               = genGroups()
	punishUser           = false
	randomToken          = envFile["RANDOM_ORG_API_TOKEN"]
)

// keyboards
var (
	noPollKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(lang["list"]),
			tgbotapi.NewKeyboardButton(lang["help"]),
		),
	)
	ownerKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(lang["list"]),
			tgbotapi.NewKeyboardButton(lang["help"]),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(lang["pollButton"]),
			tgbotapi.NewKeyboardButton(lang["shutdownButton"]),
		),
	)
	pollOwnerKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(lang["pollButton"]),
			tgbotapi.NewKeyboardButton(lang["shutdownButton"]),
		),
	)
)

// inline keyboards

var (
	noPollInline = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lang["register"], "register"),
			tgbotapi.NewInlineKeyboardButtonData(lang["quit"], "quit"),
		),
	)
	registerInline = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lang["register"], "register"),
		),
	)
	quitInline = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lang["quit"], "quit"),
		),
	)
	sendInline = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(lang["sendButton"], "send"),
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
	if randomToken == "" {
		log.Fatal("No random.org API token provided.")
	}
	if BotLanguage == "" {
		log.Fatal("No bot-language provided.")
	}
	if errParse != nil {
		log.Fatal(errParse)
	}
	// if any issues during reading the env file arise, stop the program
}
