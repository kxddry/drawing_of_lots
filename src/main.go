package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Debug = true // set true for debugging

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	startPoll := make(chan struct{})
	endPoll := make(chan struct{})
	wg := sync.WaitGroup{}
	go noPoll(updates, bot, &wg, startPoll, endPoll)
	go poll(updates, bot, &wg, startPoll, endPoll)
	wg.Add(2)
	wg.Wait()
}
