package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"sync"
)

func noPoll(c <-chan tgbotapi.Update, bot *tgbotapi.BotAPI, wg *sync.WaitGroup, startPoll chan<- struct{}, stopPoll <-chan struct{}) {
	defer wg.Done()
	for {
	loop:
		for update := range c {
			if update.Message != nil {
				chatID := update.Message.Chat.ID
				if update.Message.IsCommand() {
					if update.Message.Command() == "start" {
						text := lang["start"]
						sendNoPoll(bot, text, chatID)

						text2 := lang["start2"]
						msg := tgbotapi.NewMessage(chatID, text2)
						msg.ReplyMarkup = noPollInline
						_, _ = bot.Send(msg)
					} else {
						sendNoPoll(bot, lang["unknownCommand"], chatID)
					}
				} else {
					switch update.Message.Text {
					case lang["list"]:
						txt := ""
						switch {
						case 2 <= len(peers) && len(peers) <= 4 && (len(peers) < 10 || len(peers) > 15):
							txt = lang["listCase"] + strconv.Itoa(len(peers)) + lang["listCase1"]
						default:
							txt = lang["listCase"] + strconv.Itoa(len(peers)) + lang["listCase2"]
						}
						txt += formActiveUsers(usersHashmap)
						sendNoPoll(bot, txt, chatID)
					case lang["help"]:
						text := lang["start"]
						sendNoPoll(bot, text, chatID)

						if in(peers, chatID) != -1 {
							msg := tgbotapi.NewMessage(chatID, lang["userEnlists"])
							msg.ReplyMarkup = quitInline
							msg.ParseMode = tgbotapi.ModeHTML
							_, _ = bot.Send(msg)
						} else {
							text2 := lang["userIsThinkingAboutEnlisting"]
							msg := tgbotapi.NewMessage(chatID, text2)
							msg.ReplyMarkup = noPollInline
							_, _ = bot.Send(msg)
						}

					case lang["pollButton"]:
						if chatID == int64owner {
							txt := ""
							if len(peers) <= 1 {
								txt = lang["notEnoughUsers"]
								sendNoPoll(bot, txt, chatID)
								continue
							}
							txt = lang["pollStarted"]
							msg := tgbotapi.NewMessage(chatID, txt)
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
							err := alertMessage(bot, msg, peers)
							if err != nil {
								log.Println(err)
							}
							startPoll <- struct{}{}
							break loop
						} else {
							sendNoPoll(bot, lang["notPermitted"], chatID)
						}
					case lang["shutdownButton"]:
						if chatID == int64owner {
							sendNoPoll(bot, lang["shutdown"], int64owner)
							os.Exit(228)
						}
						sendNoPoll(bot, lang["notPermitted"], chatID)
					default:
						sendNoPoll(bot, lang["unknownMessage"], chatID)
					}
				}
			} else if update.CallbackQuery != nil {
				data := update.CallbackQuery.Data
				queryId := update.CallbackQuery.ID
				messageId := update.CallbackQuery.Message.MessageID
				userId := update.CallbackQuery.From.ID
				username := update.CallbackQuery.From.UserName
				firstname := update.CallbackQuery.From.FirstName
				switch data {
				case "register":
					txt := ""
					if in(peers, userId) == -1 {
						peers = append(peers, userId)                        // add the user to the slice of users
						usersHashmap[userId] = []string{username, firstname} // add the user to the hashmap of users
						txt = lang["addedToPoll"]
						placeholder := determinePlaceholder(userId, firstname, username)

						msg := tgbotapi.NewEditMessageText(userId, messageId, lang["userEnlists"])
						msg.ParseMode = tgbotapi.ModeHTML
						msg.ReplyMarkup = &quitInline
						_, err := bot.Send(msg)
						if err != nil {
							log.Println(err)
						}
						if userId != int64owner {
							sendNoPoll(bot, placeholder+lang["someoneAddedToPoll"], int64owner)
						}

					} else {
						txt = lang["userAlreadyInPoll"]
					}
					callback := tgbotapi.NewCallback(queryId, txt)
					_, _ = bot.Request(callback)
				case "quit":
					txt := ""
					if in(peers, userId) == -1 {
						txt = lang["userAlreadyNotInPoll"]
					} else {
						firstName := update.CallbackQuery.From.FirstName
						index := in(peers, userId)
						peers = append(peers[0:index], peers[index+1:]...)
						delete(usersHashmap, userId)
						txt = lang["userDeletedFromPoll"]

						msg := tgbotapi.NewEditMessageText(userId, messageId, lang["userDoesntEnlist"])
						msg.ParseMode = tgbotapi.ModeHTML
						msg.ReplyMarkup = &registerInline
						_, err := bot.Send(msg)
						if err != nil {
							log.Println(err)
						}

						if userId != int64owner {
							placeholder := determinePlaceholder(userId, firstName, username)
							sendNoPoll(bot, placeholder+lang["someoneDeletedFromPoll"], int64owner)
						}

					}
					callback := tgbotapi.NewCallback(queryId, txt)
					_, _ = bot.Request(callback)
				}
			}
		}
		<-stopPoll
	}
}
