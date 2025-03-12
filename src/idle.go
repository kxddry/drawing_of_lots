package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

func handleIdleBot(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, startNoPoll chan<- struct{}, startIdle <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
	loop:
		for upd := range updates {
			if upd.Message == nil {
				if upd.CallbackQuery != nil {
					chatId := upd.CallbackQuery.From.ID
					msg := tgbotapi.NewMessage(chatId, lang["idle"])
					msg.ParseMode = tgbotapi.ModeHTML
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					_, err := bot.Send(msg)
					if err != nil {
						log.Println(err)
					}
					if in(peers, chatId) == -1 {
						peers = append(peers, chatId)
					}
					usersHashmap[chatId] = []string{upd.CallbackQuery.From.UserName, upd.CallbackQuery.From.FirstName}
				}
				continue
			}
			chatId := upd.Message.Chat.ID
			if in(peers, chatId) == -1 {
				peers = append(peers, chatId)
			}
			usersHashmap[chatId] = []string{upd.Message.From.UserName, upd.Message.From.FirstName}
			if chatId != int64owner {
				msg := tgbotapi.NewMessage(chatId, lang["idle"])
				msg.ParseMode = tgbotapi.ModeHTML
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			} else {
				if upd.Message.Text == lang["startBot"] {
					msg := tgbotapi.NewMessage(upd.Message.Chat.ID, lang["starting"])
					_, _ = bot.Send(msg)
					err := alertCustom(bot, lang["start"], noPollKeyboard, ownerKeyboard, peers)
					if err != nil {
						log.Println(err)
					}
					text2 := lang["userIsThinkingAboutEnlisting"]
					err = alertCustom(bot, text2, noPollInline, noPollInline, peers)
					if err != nil {
						log.Println(err)
					}
					startNoPoll <- struct{}{}
					break loop
				} else {
					msg := tgbotapi.NewMessage(upd.Message.Chat.ID, lang["pressToStart"])
					msg.ReplyMarkup = idleOwnerKeyboard
					_, _ = bot.Send(msg)
				}
			}
		}
		<-startIdle
	}
}
