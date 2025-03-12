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

	bot.Debug = false // set true for debugging

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	startPoll := make(chan struct{})
	endPoll := make(chan struct{})
	startIdle := make(chan struct{})
	wg := sync.WaitGroup{}
	go handleIdleBot(bot, updates, endPoll, startIdle, &wg)
	go noPoll(updates, bot, &wg, startPoll, startIdle, endPoll)
	go poll(updates, bot, &wg, startPoll, endPoll, startIdle)
	wg.Add(3)
	wg.Wait()
}
