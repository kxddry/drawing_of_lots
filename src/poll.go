package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Processes the poll procedure. Runs in a goroutine automatically. The bot owner can only stop the poll.
func poll(c <-chan tgbotapi.Update, bot *tgbotapi.BotAPI,
	wg *sync.WaitGroup, startPoll <-chan struct{}, endPoll chan<- struct{}) {
	defer wg.Done()
	for {
		<-startPoll                                    // waiting for input: polling will only start when the bot owner sends /poll
		assignments := make(map[int]int64, len(peers)) // key-value pair with n --> userId
		shuffledPeers := shuffle(peers)                // pseudo randomized slice of users
		choices := make(map[int64]int, len(peers))     // store each voter's choice
		counter := Counter{}
		goodEnd := false
		myUpdChan := make(chan tgbotapi.Update)
		go counter.process(myUpdChan, bot, choices)

		for n, participant := range shuffledPeers {
			assignments[n] = participant
		}

		table := "Очередь:\n" + formTable(usersHashmap, assignments)
		for n, chatID := range assignments {
			txt := "Вам был назначен номер " + strconv.Itoa(n) + ". Вы " + strconv.Itoa(n+1) + " в порядке очереди."
			msg := tgbotapi.NewMessage(chatID, txt)
			if chatID == int64owner {
				msg.ReplyMarkup = pollOwnerKeyboard
			}
			_, _ = bot.Send(msg)
			_ = send(bot, table, chatID)
		}

		i := 0

		question := "Выбор группы"
		options := genGroups()
		thePoll := tgbotapi.NewPoll(assignments[i], question, options...)
		thePoll.IsAnonymous = false

		pollText := "Опрос предоставлен ниже. \n\n" + "Когда сделаете окончательный выбор, пожалуйста, нажмите на кнопку ниже.\n\n"

		pollMsg := tgbotapi.NewMessage(assignments[i], pollText)
		_, _ = bot.Send(pollMsg)

		_ = sendCounter(bot, &counter, assignments[i])
		msg, err := bot.Send(thePoll)
		if err != nil {
			log.Fatal(err)
		}
		button := tgbotapi.NewMessage(assignments[i], "кнопка ↓")
		button.ParseMode = tgbotapi.ModeHTML
		button.ReplyMarkup = sendInline
		_, _ = bot.Send(button)
		tellQueue := func() {
			txt := "Сейчас выбирает " +
				determinePlaceholder(assignments[i], usersHashmap[assignments[i]][1], usersHashmap[assignments[i]][0])
			_ = alertEveryoneButTXT(assignments[i], bot, txt, peers)
		}
		tellQueue()
		for update := range c {
			if update.PollAnswer != nil && update.PollAnswer.PollID != msg.Poll.ID {
				continue
			}

			if update.PollAnswer != nil && update.PollAnswer.User.ID == assignments[i] {
				myUpdChan <- update
			} else if update.PollAnswer != nil && update.PollAnswer.User.ID != assignments[i] {
				chatID := update.PollAnswer.User.ID
				txt := determinePlaceholder(chatID, usersHashmap[chatID][1], usersHashmap[chatID][0])
				txt += " повлиял на голосование <b>не в свою очередь</b> ⚠️⚠️⚠️ \n\nТеперь придётся проводить голосование заново."
				_ = alertEveryoneButTXT(chatID, bot, txt, peers)
				_ = send(bot, "Вы проголосовали не в свою очередь. Теперь придётся провести голосование заново.", chatID)
				break
			}
			if update.Message != nil {
				chatID := update.Message.From.ID
				if chatID != int64owner {
					// only process messages
					// from the owner and the curr id
					_ = send(bot, "Сообщения в режиме голосования не принимаются.", chatID)
					continue
				}

				if chatID == int64owner {
					if update.Message.Command() == "poll" {
						err := send(bot, "poll stopped", int64owner)
						if err != nil {
							log.Println(err)
						}
						break
					} else if update.Message.Command() == "shutdown" {
						err := send(bot, "Shutting down.", int64owner)
						if err != nil {
							log.Fatal(err)
						}
						os.Exit(0)
					}
				}
			} else if update.CallbackQuery != nil {
				chatID := update.CallbackQuery.From.ID
				query := update.CallbackQuery.Data
				queryID := update.CallbackQuery.ID
				messageID := update.CallbackQuery.Message.MessageID
				if query == "send" {
					if chatID == assignments[i] {

						if choices[chatID] == -1 {
							_, _ = bot.Request(tgbotapi.NewCallback(queryID, "Вы должны выбрать хоть какую-нибудь группу."))
							continue
						}

						if punishUser {
							_, _ = bot.Request(tgbotapi.NewCallback(queryID, "Вы должны выбрать группу, в которой есть места."))
						}

						i++

						txt := "Вы выбрали " + strings.Replace(groups[choices[chatID]], "ппа", "ппу", -1) + "."
						editor := tgbotapi.NewEditMessageText(chatID, messageID, "Спасибо за ответ. "+txt)
						_, _ = bot.Send(editor)

						forward := tgbotapi.NewForward(assignments[i], assignments[i-1], msg.MessageID)
						if i == len(peers) { // end the poll
							// forward the poll to the owner
							forward = tgbotapi.NewForward(int64owner, assignments[i-1], msg.MessageID)
							txt = "Голосование завершено. Всем спасибо за участие. " +
								"Опрос отправлен " + determinePlaceholder(int64owner,
								usersHashmap[int64owner][1], usersHashmap[int64owner][0])
							_ = alertCustomBut(int64owner, bot, txt, noPollKeyboard, ownerKeyboard, peers)
							_, _ = bot.Send(forward)

							// delete the poll from the chat with the last guy
							deleter := tgbotapi.NewDeleteMessage(assignments[i-1], msg.MessageID)
							_, _ = bot.Send(deleter)
							_ = send(bot, "Голосование завершено. Результаты выше.", int64owner)
							goodEnd = true
							break
						}
						// don't end the poll, then
						tellQueue()

						pollMsg = tgbotapi.NewMessage(assignments[i], pollText)
						_, _ = bot.Send(pollMsg)
						infoText := "Выбирайте группу, в которой есть места. Счётчик мест для каждой группы: \n\n" +
							formCounter(&counter) + "\n\n" + "Подходят группы, в которых <b>меньше " +
							strconv.Itoa(MaxUsersPerGroup) + "</b> человек. Люди, выбравшие группы:\n\n" + formChosen(choices)
						_ = send(bot, infoText, assignments[i])
						tmp, _ := bot.Send(forward)
						deleter := tgbotapi.NewDeleteMessage(assignments[i-1], msg.MessageID)
						_, _ = bot.Send(deleter)
						msg = tmp

						button = tgbotapi.NewMessage(assignments[i], "кнопка ↓")
						button.ParseMode = tgbotapi.ModeHTML
						button.ReplyMarkup = sendInline
						_, _ = bot.Send(button)

						callback := tgbotapi.NewCallback(queryID, "Спасибо за уделённое время.")
						_, _ = bot.Request(callback)

					} else { // chatID != assignments[i]
						callback := tgbotapi.NewCallback(queryID, "забанят")
						_, _ = bot.Request(callback)
					}
				} else { // non-send query (probably noPoll query)
					callback := tgbotapi.NewCallback(queryID, "Дождитесь окончания голосования.")
					_, _ = bot.Request(callback)
				}
			}
		}

		if !goodEnd {
			txt := "Что-то пошло не так во время голосования. Голосование завершено.\n\n" +
				"Зарегистрируйтесь на новое голосование, нажав кнопку ниже."
			_ = alertCustom(bot, txt, noPollKeyboard, ownerKeyboard, peers)
		}
		// end of the poll -- reset all data
		peers = make([]int64, 0, len(peers))
		shuffledPeers = make([]int64, 0, len(peers))
		usersHashmap = make(map[int64][]string, len(usersHashmap))
		assignments = make(map[int]int64, len(assignments))

		endPoll <- struct{}{}

	}
}
