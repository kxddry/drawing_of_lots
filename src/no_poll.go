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
						text := "Добро пожаловать! \n<b>Использование:</b> \n\t- \"<b>Регистрация</b>\" для регистрации в голосовании" +
							" \n\t- \"<b>Выход</b>\", чтобы не участвовать" +
							"\n\t- \"<b>Помощь</b>\" для помощи \n\t- \"<b>Список</b>\" для получения списка участвующих."
						sendNoPoll(bot, text, chatID)

						text2 := "Выберите, будете ли участвовать в голосовании."
						msg := tgbotapi.NewMessage(chatID, text2)
						msg.ReplyMarkup = noPollInline
						_, _ = bot.Send(msg)
					} else {
						sendNoPoll(bot, "Неизвестная команда.", chatID)
					}
				} else {
					switch update.Message.Text {
					case "Список":
						txt := ""
						switch {
						case 2 <= len(peers) && len(peers) <= 4 && (len(peers) < 10 || len(peers) > 15):
							txt = "Сейчас участвует " + strconv.Itoa(len(peers)) + " человека в голосовании.\n"
						default:
							txt = "Сейчас участвует " + strconv.Itoa(len(peers)) + " человек в голосовании.\n"
						}
						txt += formActiveUsers(usersHashmap)
						sendNoPoll(bot, txt, chatID)
					case "Помощь":
						text := "Добро пожаловать! \n<b>Использование:</b> \n\t- \"<b>Регистрация</b>\" для регистрации в голосовании" +
							" \n\t- \"<b>Выход</b>\", чтобы не участвовать" +
							"\n\t- \"<b>Помощь</b>\" для помощи \n\t- \"<b>Список</b>\" для получения списка участвующих."
						sendNoPoll(bot, text, chatID)

						text2 := "Выберите, будете ли участвовать в голосовании."
						msg := tgbotapi.NewMessage(chatID, text2)
						msg.ReplyMarkup = noPollInline
						_, _ = bot.Send(msg)

					case "Poll":
						if chatID == int64owner {
							txt := ""
							if len(peers) <= 1 {
								txt = "not enough users"
								err := send(bot, txt, chatID)
								if err != nil {
									log.Println(err)
								}
								continue
							}
							txt = "Голосование началось."
							msg := tgbotapi.NewMessage(chatID, txt)
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
							err := alertMessage(bot, msg, peers)
							if err != nil {
								log.Println(err)
							}
							startPoll <- struct{}{}
							break loop
						} else {
							_ = send(bot, "You're not permitted to do that.", chatID)
						}
					case "Shutdown":
						if chatID == int64owner {
							err := send(bot, "Shutting down.", int64owner)
							if err != nil {
								log.Fatal(err)
							}
							os.Exit(228)
						}
						_ = send(bot, "You're not permitted to do that.", chatID)
					default:
						_ = send(bot, "Я не знаю, как ответить на это сообщение.", chatID)
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
						txt = "Вы добавлены в голосование."
						placeholder := determinePlaceholder(userId, firstname, username)

						msg := tgbotapi.NewEditMessageText(userId, messageId, "Вы <b>участвуете</b> в голосовании.")
						msg.ParseMode = tgbotapi.ModeHTML
						msg.ReplyMarkup = &quitInline
						_, err := bot.Send(msg)
						if err != nil {
							log.Println(err)
						}
						if userId != int64owner {
							err := send(bot, placeholder+" добавлен(-а) в голосование.", int64owner)
							if err != nil {
								log.Println(err)
							}
						}

					} else {
						txt = "Вы уже участвуете в голосовании."
					}
					callback := tgbotapi.NewCallback(queryId, txt)
					_, _ = bot.Request(callback)
				case "quit":
					txt := ""
					if in(peers, userId) == -1 {
						txt = "Вы и так не участвуете."
					} else {
						firstName := update.CallbackQuery.From.FirstName
						index := in(peers, userId)
						peers = append(peers[0:index], peers[index+1:]...)
						delete(usersHashmap, userId)
						txt = "Вы удалены из голосования."

						msg := tgbotapi.NewEditMessageText(userId, messageId, "Вы <b>не</b> участвуете в голосовании.")
						msg.ParseMode = tgbotapi.ModeHTML
						msg.ReplyMarkup = &registerInline
						_, err := bot.Send(msg)
						if err != nil {
							log.Println(err)
						}

						if userId != int64owner {
							placeholder := determinePlaceholder(userId, firstName, username)
							err := send(bot, placeholder+" удален(-а) из голосования.", int64owner)
							if err != nil {
								log.Println(err)
							}
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
