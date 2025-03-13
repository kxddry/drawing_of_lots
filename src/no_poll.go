package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"sync"
)

func noPoll(c <-chan tgbotapi.Update, bot *tgbotapi.BotAPI, wg *sync.WaitGroup, startPoll, startIdle chan<- struct{}, stopPoll <-chan struct{}) {
	defer wg.Done()
	for {
		<-stopPoll
	loop:
		for update := range c {
			if update.Message != nil {
				chatID := update.Message.Chat.ID
				firstName := update.Message.From.FirstName
				username := update.Message.From.UserName
				updateDatabase(chatID, username, firstName)
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
						length := func() int {
							count := 0
							for _, v := range participants {
								if v == 1 {
									count++
								}
							}
							return count
						}()
						switch {
						case 2 <= length && length <= 4 && (length < 10 || length > 15):
							txt = lang["listCase"] + strconv.Itoa(length) + lang["listCase1"]
						default:
							txt = lang["listCase"] + strconv.Itoa(length) + lang["listCase2"]
						}
						txt += formActiveUsers()
						sendNoPoll(bot, txt, chatID)
					case lang["help"]:
						text := lang["start"]
						sendNoPoll(bot, text, chatID)

						if participants[chatID] == 1 {
							msg := tgbotapi.NewMessage(chatID, lang["userEnlists"])
							msg.ReplyMarkup = quitInline
							msg.ParseMode = tgbotapi.ModeHTML
							_, _ = bot.Send(msg)
						} else if participants[chatID] == 0 {
							text2 := lang["userIsThinkingAboutEnlisting"]
							msg := tgbotapi.NewMessage(chatID, text2)
							msg.ReplyMarkup = noPollInline
							msg.ParseMode = tgbotapi.ModeHTML
							_, _ = bot.Send(msg)
						} else {
							text2 := lang["userDoesntEnlist"]
							msg := tgbotapi.NewMessage(chatID, text2)
							msg.ReplyMarkup = registerInline
							msg.ParseMode = tgbotapi.ModeHTML
							_, _ = bot.Send(msg)
						}

					case lang["pollButton"]:
						if chatID == int64owner {
							txt := ""
							length := func() int {
								count := 0
								for _, v := range participants {
									if v == 1 {
										count++
									}
								}
								return count
							}()
							if length <= 1 {
								txt = lang["notEnoughUsers"]
								sendNoPoll(bot, txt, chatID)
								continue
							}
							txt = lang["pollStarted"]
							msg := tgbotapi.NewMessage(chatID, txt)
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
							err := alertMessage(bot, msg, getActivePeers())
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
							err := alertCustom(bot, lang["shutdown"], tgbotapi.NewRemoveKeyboard(true), idleOwnerKeyboard, peers)
							if err != nil {
								log.Println(err)
							}
							participants = make(map[int64]int, len(participants))
							startIdle <- struct{}{}
							break loop
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
				updateDatabase(userId, username, firstname)
				switch data {
				case "register":
					txt := ""
					if participants[userId] <= 0 {
						participants[userId] = 1
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
					if participants[userId] == -1 {
						txt = lang["userAlreadyNotInPoll"]
					} else {
						firstName := update.CallbackQuery.From.FirstName
						participants[userId] = -1
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
	}
}
