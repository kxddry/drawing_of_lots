package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"sync"
)

func noPoll(c <-chan tgbotapi.Update, bot *tgbotapi.BotAPI, wg *sync.WaitGroup, startPoll chan<- struct{}, stopPoll <-chan struct{}) {
	defer wg.Done()
	for {
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
				txt = "Добро пожаловать! \nИспользование: \n\t- /register для регистрации в голосовании" +
					" \n\t- /quit, чтобы не участвовать" +
					"\n\t- /help для помощи \n\t- /list для получения списка участвующих."
			case "register":
				if in(peers, chatID) == -1 {
					username := update.Message.From.UserName
					firstName := update.Message.From.FirstName

					peers = append(peers, chatID)                        // add the user to the slice of users
					usersHashmap[chatID] = []string{username, firstName} // add the user to the hashmap of users
					txt = "Вы добавлены в голосование."

					if chatID != int64owner {
						err := send(bot, firstName+" добавлен(-а) в голосование.", int64owner)
						if err != nil {
							log.Println(err)
						}
					}
				} else {
					txt = "Вы уже участвуете в голосовании."
				}
			case "quit":
				if in(peers, chatID) == -1 {
					txt = "Вы и так не участвуете."
				} else {
					firstName := update.Message.From.FirstName
					index := in(peers, chatID)
					peers = append(peers[0:index], peers[index+1:]...)
					delete(usersHashmap, chatID)
					txt = "Вы удалены из голосования."

					if chatID != int64owner {
						err := send(bot, firstName+" удален(-а) из голосования.", int64owner)
						if err != nil {
							log.Println(err)
						}
					}
				}
			case "help":
				txt = "Использование: \n\t- /register для регистрации в голосовании \n\t- /quit, чтобы не участвовать" +
					"\n\t- /help для помощи \n\t- /list для получения списка участвующих."
			case "list":
				switch {
				case 2 <= len(peers) && len(peers) <= 4 && (len(peers) < 10 || len(peers) > 15):
					txt = "Сейчас участвует " + strconv.Itoa(len(peers)) + " человека в голосовании.\n"
				default:
					txt = "Сейчас участвует " + strconv.Itoa(len(peers)) + " человек в голосовании.\n"
				}
				txt += formActiveUsers(usersHashmap)
			default:
				txt = "Неизвестная команда."
			}

			if chatID == int64owner {
				if update.Message.Command() == "poll" {
					if len(peers) <= 1 {
						txt = "not enough users"
						err := send(bot, txt, chatID)
						if err != nil {
							log.Println(err)
						}
						continue
					}
					txt = "Голосование началось."
					err := alertEveryone(bot, txt, peers)
					if err != nil {
						log.Println(err)
					}
					startPoll <- struct{}{}
					break
				} else if update.Message.Command() == "list" {
					txt = "There are currently " + strconv.Itoa(len(peers)) + " people in the poll. \n"
					txt += formActiveUsers(usersHashmap)
				}

			}
			if txt != "" {
				err := send(bot, txt, chatID)
				if err != nil {
					log.Println(err)
				}
			}
		}
		<-stopPoll
	}
}
