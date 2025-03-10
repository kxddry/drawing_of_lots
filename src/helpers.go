package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strconv"
)

func in(peers []int64, target int64) int {
	if len(peers) == 0 {
		return -1
	}
	for index, peer := range peers {
		if peer == target {
			return index
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

func formTable(nicknamesAndIDs map[int64][]string, assignments map[int]int64) string {
	res := ""
	for n, chatID := range assignments {
		username, firstName := nicknamesAndIDs[chatID][0], nicknamesAndIDs[chatID][1]
		str := ""
		if username == "" {
			str = strconv.Itoa(n) + " --> " + firstName + "\n"
		} else {
			str = strconv.Itoa(n) + " --> @" + username + " (" + firstName + ")\n"
		}
		res += str
	}
	return res
}

func formActiveUsers(nicknamesAndIDs map[int64][]string) string {
	res := ""
	i := 1
	for _, arr := range nicknamesAndIDs {
		username, firstName := arr[0], arr[1]
		str := ""
		if username == "" {
			str = strconv.Itoa(i) + ". " + firstName + "\n"
		} else {
			str = strconv.Itoa(i) + ". @" + username + " (" + firstName + ")\n"
		}
		res += str
		i++
	}
	return res
}

func send[sendable string | []byte | []rune](bot *tgbotapi.BotAPI, txt sendable, chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, string(txt))
	_, err := bot.Send(msg)
	return err
}

func alertEveryone[sendable string | []byte](bot *tgbotapi.BotAPI, txt sendable, peers []int64) error {
	for _, peer := range peers {
		err := send(bot, txt, peer)
		if err != nil {
			return err
		}
	}
	return nil
}
