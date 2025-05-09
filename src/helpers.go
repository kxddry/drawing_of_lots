package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
				_ = alertEveryoneButTXT(id, bot, placeholder+lang["cancelledTheirVote"], getActivePeers())
			} else {
				punishUser = false
			}
			if choices[id] == -1 {
				_ = send(bot, lang["handleBug"], id)
			}
			c.arr[choices[id]]--
			choices[id] = -1
		} else {
			choice := upd.PollAnswer.OptionIDs[0]
			choices[id] = choice
			c.arr[choice]++
			if c.arr[choice] > MaxUsersPerGroup {
				punishUser = true
				_ = send(bot, lang["groupFilled"], id)
				_ = alertEveryoneButTXT(id, bot, "⚠️⚠️⚠️ "+placeholder+lang["someoneTriedFilledGroup"]+
					"⚠️⚠️⚠️", getActivePeers())
			} else {
				if BotLanguage == "russian" {
					group := strings.Replace(groups[choice], "ппа", "ппу", -1) // russian language workarounds
					_ = alertEveryoneButTXT(id, bot, placeholder+lang["chose"]+group+".", getActivePeers())
				} else {
					_ = alertEveryoneButTXT(id, bot, placeholder+lang["chose"]+groups[choice]+".", getActivePeers())

				}
			}
		}
	}
}

func genGroups() []string {
	res := make([]string, 0, NumberOfGroups)
	for i := 1; i != NumberOfGroups+1; i++ {
		res = append(res, lang["group"]+strconv.Itoa(i))
	}
	return res
}

func formTable(assignments map[int]int64) string {
	res := ""
	for n := 0; n != len(assignments); n++ {
		chatID := assignments[n]
		username, firstName := usersHashmap[chatID][0], usersHashmap[chatID][1]
		str := strconv.Itoa(n) + " --> " + determinePlaceholder(chatID, firstName, username) + "\n"
		res += str
	}
	return res
}

func formActiveUsers() string {
	res := ""
	i := 1
	for ID, arr := range usersHashmap {
		if participants[ID] < 1 {
			continue
		}
		username, firstName := arr[0], arr[1]
		str := strconv.Itoa(i) + ". " + determinePlaceholder(ID, firstName, username) + "\n"
		res += str
		i++
	}
	return res
}

func send[sendable string | []byte | []rune](bot *tgbotapi.BotAPI, txt sendable, chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, string(txt))
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.Send(msg)
	return err
}

func alertMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig, peers []int64) error {
	for _, peerId := range peers {
		msg.BaseChat.ChatID = peerId
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
-----UNUSED FUNCTIONS-----

	func alertEveryoneTXT[sendable string | []byte](bot *tgbotapi.BotAPI, txt sendable, peers []int64) error {
		for _, peer := range peers {
			err := send(bot, txt, peer)
			if err != nil {
				return err
			}
		}
		return nil
	}

	func alertMsgBut(id int64, bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig, peers []int64) error {
		for _, peer := range peers {
			if peer != id {
				msg.BaseChat.ChatID = peer
				_, err := bot.Send(msg)
				if err != nil {
					return err
				}
			}
		}
		return nil
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

-----END UNUSED FUNCTIONS-----
*/
func alertCustomBut(id int64, bot *tgbotapi.BotAPI, txt string, reply, replyOwner tgbotapi.ReplyKeyboardMarkup, peers []int64) error {
	for _, peer := range peers {
		if peer != id {
			msg := tgbotapi.NewMessage(peer, txt)
			msg.ParseMode = tgbotapi.ModeHTML
			msg.BaseChat.ChatID = peer
			if peer == int64owner {
				msg.ReplyMarkup = replyOwner
			} else {
				msg.ReplyMarkup = reply
			}
			_, err := bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func alertCustom(bot *tgbotapi.BotAPI,
	txt string, reply, replyOwner interface{}, peers []int64) error {
	for _, peer := range peers {
		msg := tgbotapi.NewMessage(peer, txt)
		msg.BaseChat.ChatID = peer
		msg.ParseMode = tgbotapi.ModeHTML
		if peer == int64owner {
			msg.ReplyMarkup = replyOwner
		} else {
			msg.ReplyMarkup = reply
		}
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func alertEveryoneButTXT[sendable string | []byte](id int64, bot *tgbotapi.BotAPI, txt sendable, peers []int64) error {
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
		res += "- <b>" + groups[i] + ": " + strconv.Itoa(c.arr[i]) + "</b>\n"
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
	msg.ParseMode = tgbotapi.ModeHTML
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

func sendNoPoll(bot *tgbotapi.BotAPI, txt string, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, txt)
	if chatID == int64owner {
		msg.ReplyMarkup = ownerKeyboard
	} else {
		msg.ReplyMarkup = noPollKeyboard
	}
	msg.ParseMode = tgbotapi.ModeHTML
	_, _ = bot.Send(msg)
}

func initChoices(choices *map[int64]int, peers []int64) {
	for _, peer := range peers {
		(*choices)[peer] = -1
	}
}

func getActivePeers() []int64 {
	res := make([]int64, 0, cap(peers))
	for _, id := range peers {
		if participants[id] == 1 {
			res = append(res, id)
		}
	}
	return res
}

func updateDatabase(chatId int64, username, firstName string) {
	if len(usersHashmap[chatId]) == 0 {
		usersHashmap[chatId] = []string{username, firstName}
		peers = append(peers, chatId)
	}
}
