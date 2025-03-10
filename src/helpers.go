package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
)

type Counter struct {
	arr [NumberOfGroups]int
}

func (c *Counter) process(updChan <-chan tgbotapi.Update, bot *tgbotapi.BotAPI, choices map[int64]int) {
	for upd := range updChan {
		id := upd.PollAnswer.User.ID
		username := upd.PollAnswer.User.UserName
		firstname := upd.PollAnswer.User.FirstName
		placeholder := determinePlaceholder(id, firstname, username)
		if len(upd.PollAnswer.OptionIDs) == 0 {
			if c.arr[choices[id]] <= 5 {
				_ = alertEveryoneBut(id, bot, placeholder+" отменил(-а) свой голос!", peers)
			} else {
				punishUser = false
			}
			if choices[id] == -1 {
				_ = send(bot, "забанят", id)
			}
			c.arr[choices[id]]--
			choices[id] = -1
		} else {
			choice := upd.PollAnswer.OptionIDs[0]
			choices[id] = choice
			c.arr[choice]++
			if c.arr[choice] > MaxUsersPerGroup {
				punishUser = true
				_ = send(bot, "Вам нужно перевыбрать. Эта группа заполнена доверху.", id)
				_ = alertEveryoneBut(id, bot, "⚠️⚠️⚠️ "+placeholder+" попытался(-ась) выбрать "+
					"заполненную группу ⚠️⚠️⚠️", peers)
			} else {
				group := strings.Replace(groups[choice], "ппа", "ппу", -1) // russian language workarounds
				_ = alertEveryoneBut(id, bot, placeholder+" выбрал(-а) "+group+".", peers)
			}
		}
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
	for n := 0; n != len(assignments); n++ {
		chatID := assignments[n]
		username, firstName := nicknamesAndIDs[chatID][0], nicknamesAndIDs[chatID][1]
		str := strconv.Itoa(n) + " --> " + determinePlaceholder(chatID, firstName, username) + "\n"
		res += str
	}
	return res
}

func formActiveUsers(nicknamesAndIDs map[int64][]string) string {
	res := ""
	i := 1
	for ID, arr := range nicknamesAndIDs {
		username, firstName := arr[0], arr[1]
		str := strconv.Itoa(i) + ". " + determinePlaceholder(ID, firstName, username) + "\n"
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

func formCounter(c *Counter) string {
	res := ""
	allZeroes := true
	for i := 0; i < NumberOfGroups; i++ {
		if c.arr[i] != 0 {
			allZeroes = false
			break
		}
	}
	if allZeroes {
		return res
	}
	for i := 0; i < NumberOfGroups; i++ {
		res += "- " + groups[i] + ": " + strconv.Itoa(c.arr[i]) + "\n"
	}
	return res
}

func determinePlaceholder(id int64, firstname, username string) string {
	placeholder := ""
	if username == "" { // empty username
		if !checkValid(firstname) { // empty username, non-valid first name
			placeholder = strconv.Itoa(int(id))
		} else {
			placeholder = firstname
		}
	} else {
		placeholder = "@" + username
	}
	return placeholder
}

func sendCounter(bot *tgbotapi.BotAPI, c *Counter, chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, formCounter(c))
	_, err := bot.Send(msg)
	return err
}

func formChosen(choices map[int64]int) string {
	res := ""
	for i := 0; i < NumberOfGroups; i++ {
		str := "- " + groups[i] + ": [ "
		for userId, choice := range choices {
			username, firstName := usersHashmap[userId][0], usersHashmap[userId][1]
			placeholder := determinePlaceholder(userId, firstName, username)
			if choice == i {
				str += placeholder + " "
			}
		}
		str += "]\n"
		res += str
	}
	return res
}
