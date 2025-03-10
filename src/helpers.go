package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strconv"
	"unicode"
)

type Counter struct {
	arr [NumberOfGroups]int
}

func (c *Counter) process(upd tgbotapi.Update, bot *tgbotapi.BotAPI, choices map[int64]int) {
	username := upd.Message.From.UserName
	firstname := upd.Message.From.FirstName
	id := upd.Message.From.ID
	placeholder := ""
	choice := upd.PollAnswer.OptionIDs[0]
	if username == "" { // empty username
		if !checkValid(firstname) { // empty username, non-valid first name
			placeholder = strconv.Itoa(int(id))
		}
		placeholder = firstname
	} else {
		placeholder = username
	}
	if len(upd.PollAnswer.OptionIDs) == 0 {
		_ = alertEveryoneBut(id, bot, placeholder+" отменил(-а) свой голос!", peers)
	} else {
		_ = alertEveryoneBut(id, bot, placeholder+" выбрал "+groups[choice], peers)
		choices[id] = choice
		c.arr[choice]++
	}
}

func genGroups() []string {
	res := make([]string, 0, NumberOfGroups)
	for i := 1; i != NumberOfGroups+1; i++ {
		res = append(res, "группа "+strconv.Itoa(i))
	}
	return res
}

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

func alertEveryoneBut[sendable string | []byte](id int64, bot *tgbotapi.BotAPI, txt sendable, peers []int64) error {
	for _, peer := range peers {
		if peer != id {
			err := send(bot, txt, peer)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkValid(str string) bool {
	for _, char := range str {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			return true
		}
	}
	return false
}
