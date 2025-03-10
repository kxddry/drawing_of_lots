package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"sync"
)

func poll(c <-chan tgbotapi.Update, bot *tgbotapi.BotAPI,
	wg *sync.WaitGroup, startPoll <-chan struct{}, endPoll chan<- struct{}) {
	defer wg.Done()
	for {
		<-startPoll
		assignments := make(map[int]int64, len(peers))
		shuffledPeers := shuffle(peers)

		for n, participant := range shuffledPeers {
			assignments[n] = participant
		}

		table := "Очередь:\n" + formTable(usersHashmap, assignments)
		for n, chatID := range assignments {
			txt := "Вам был назначен номер " + strconv.Itoa(n) + ". Вы " + strconv.Itoa(n+1) + " в порядке очереди."

			err := send(bot, txt, chatID)
			err2 := send(bot, table, chatID)
			if err != nil || err2 != nil {
				log.Println(err, err2)
			}
		}

		i := 0

		question := "Выбор группы"
		options := []string{"группа 1", "группа 2", "группа 3"}
		thePoll := tgbotapi.NewPoll(assignments[i], question, options...)
		thePoll.IsAnonymous = false

		pollText := "Вот опрос. \n\n" + "Когда сделаете окончательный выбор, пожалуйста, нажмите /next.\n\n"
		err := send(bot, pollText, assignments[i])
		if err != nil {
			log.Println(err)
		}
		msg, err := bot.Send(thePoll)
		if err != nil {
			log.Fatal(err)
		}

		for update := range c {
			if update.Message == nil {
				continue
			}
			chatID := update.Message.Chat.ID
			if chatID != int64owner && chatID != assignments[i] {
				// only process messages
				// from the owner and the curr id
				continue
			}
			if chatID == assignments[i] && update.Message.Command() == "next" {
				i++
				_ = send(bot, "Вы сделали свой выбор.", chatID)
				forward := tgbotapi.NewForward(assignments[i], assignments[i-1], msg.MessageID)
				if i == len(peers) {
					forward = tgbotapi.NewForward(int64owner, assignments[i-1], msg.MessageID)
					txt := "Голосование завершено. Всем спасибо за участие. " +
						"Опрос отправлен @" + usersHashmap[int64owner][0]
					_, _ = bot.Send(forward)
					deleter := tgbotapi.NewDeleteMessage(assignments[i-1], msg.MessageID)
					_, _ = bot.Send(deleter)
					_ = alertEveryone(bot, txt, peers)
					break
				}
				_ = send(bot, pollText, assignments[i])
				tmp, _ := bot.Send(forward)
				deleter := tgbotapi.NewDeleteMessage(assignments[i-1], msg.MessageID)
				_, _ = bot.Send(deleter)
				msg = tmp
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
		}

		peers = make([]int64, 0, len(peers))
		assignments = make(map[int]int64, len(assignments))
		endPoll <- struct{}{}
		shuffledPeers = make([]int64, 0, len(peers))
		usersHashmap = make(map[int64][]string, len(usersHashmap))

	}
}
