package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

func send[sendable string | []byte | []rune](bot *tgbotapi.BotAPI, txt sendable, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, string(txt))
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func no_poll(c <-chan tgbotapi.Update, bot *tgbotapi.BotAPI, wg *sync.WaitGroup, ch2 chan<- struct{}) {
	defer wg.Done()
	for update := range c {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		if !update.Message.IsCommand() {
			continue
		}

		txt := ""
		switch update.Message.Command() {
		case "start":
			if in(peers, chatID) == -1 {
				peers = append(peers, chatID)
				txt = "Вы добавлены в голосование."
			} else {
				txt = "Вы уже участвуете в голосовании."
			}
		case "quit":
			if in(peers, chatID) == -1 {
				txt = "Вы и так не участвуете."
			} else {
				index := in(peers, chatID)
				peers = append(peers[0:index], peers[index+1:]...)
			}
		}

		if chatID == int64owner {
			switch update.Message.Command() {
			case "poll":
				ch2 <- struct{}{}
				wg.Done()
				return
				// TODO: call the /poll function to start the poll
			}
		}
		send(bot, txt, chatID)
	}
}
