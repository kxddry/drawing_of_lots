package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"strconv"
	"sync"
)

var (
	envFile, errOpen = godotenv.Read("pkg/.env")
	token            = envFile["TOKEN"]
	owner            = envFile["OWNER"]
	int64owner, _    = strconv.ParseInt(owner, 10, 64)
	peers            = make([]int64, 30)
)

func init() {
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	if token == "" {
		log.Fatal("check your API token")
	}
	if owner == "" {
		log.Fatal("No owner provided.")
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	start_poll := make(chan struct{})
	end_poll := make(chan struct{})
	wg := sync.WaitGroup{}
	go no_poll(updates, bot, &wg, start_poll)
	go poll(updates, bot, &wg, start_poll, end_poll)
	wg.Add(2)
	wg.Wait()
}

// TODO: this function should probably also take map as an argument
func poll(c <-chan tgbotapi.Update, bot *tgbotapi.BotAPI,
	wg *sync.WaitGroup, signal <-chan struct{}, end chan<- struct{}) {
	<-signal
	defer wg.Done()
	// TODO: manage polls
	assignments := make(map[int]int64, len(peers))
	shuffled_peers := shuffle(peers)
	for i, v := range shuffled_peers {
		assignments[i] = v
	}
	for i, v := range assignments {
		txt := "Вам был назначен номер " + strconv.Itoa(i) + ". Вы " + strconv.Itoa(i+1) + " в порядке очереди."
		send(bot, txt, v)
	}
	i := 0
	for update := range c {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		if chatID != int64owner || (!update.Message.IsCommand() && update.Message.Poll == nil) {
			continue
			// if the message we're receiving is not from the owner
			// NOR is it a command
			// NOR is it a poll, then we ignore
		}
		pollID := -1
		if update.Message.Poll != nil && chatID == int64owner {
			pollID = update.Message.MessageID
		}
		switch update.Message.Command() {
		case "send":
			if pollID == -1 {
				send(bot, "send poll first", int64owner)
				continue
			}
			msg := tgbotapi.NewForward(assignments[i], int64owner, pollID)
			_, err := bot.Send(msg)
			if err != nil {
				send(bot, fmt.Sprintf("%v", err), int64owner)
			}
		case "next":
			i++
		case "end":
			end <- struct{}{}
			return
		}

	}

}

// msg.ReplyToMessageID = update.Message.MessageID  // TODO: если нужно ответить на сообщение

func in(peers []int64, target int64) int {
	if len(peers) == 0 {
		return -1
	}
	for i, v := range peers {
		if v == target {
			return i
		}
	}
	return -1
}

func shuffle(slice []int64) []int64 {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
