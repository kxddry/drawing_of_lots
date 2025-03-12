package main

// language: message hashmaps
// and I know JSON is probably a better way of storing data, but I'm too lazy to marshal
// and unmarshal data, create structures, and so on. So this is probably easier
// as it's a small project anyway.
var (
	language = map[string]map[string]string{
		"english": {
			// idle
			"startBot":     "Start bot",
			"pressToStart": "Press the button below to start the bot.",
			"startingBot":  "Okay, starting the bot...",
			"idle":         "This bot is currently disabled. Ask the owner to start the poll.",

			// no poll
			"register": "Register",
			"quit":     "Quit",
			"start": "Welcome! \n<b>Usage:</b> \n\t- \"<b>Register</b>\" to register in the poll" +
				" \n\t- \"<b>Quit</b>\" to quit the poll" +
				"\n\t- \"<b>Help</b>\" to display help \n\t- \"<b>List</b>\" to get the list of active poll participants.",
			"start2":                       "Choose whether you're partaking in the poll.",
			"unknownCommand":               "Unknown command",
			"list":                         "List",
			"help":                         "Help",
			"listCase":                     "There are ",
			"listCase1":                    " people in the poll right now.\n",
			"listCase2":                    " people in the poll right now.\n",
			"userEnlists":                  "You <b>are taking</b> the poll.",
			"userAlreadyInPoll":            "You are already registered.",
			"userAlreadyNotInPoll":         "You are not taking the poll already.",
			"userDeletedFromPoll":          "You are removed from the poll.",
			"userDoesntEnlist":             "You <b>are not</b> taking the poll.",
			"userIsThinkingAboutEnlisting": "Choose whether you're partaking in the poll.",
			"notEnoughUsers":               "Not enough users.",
			"pollStarted":                  "The poll has started.",
			"notPermitted":                 "You are not permitted to do that action.",
			"shutdown":                     "The bot is now disabled.",
			"unknownMessage":               "I do not know how to respond to that.",
			"addedToPoll":                  "You are added to the poll.",
			"someoneAddedToPoll":           " added to the poll.",
			"someoneDeletedFromPoll":       " removed from the poll.",

			// helpers

			"cancelledTheirVote":      " cancelled their vote!",
			"handleBug":               "you're going to get banned",
			"groupFilled":             "You have to repick. This group is filled to the brim.",
			"someoneTriedFilledGroup": " tried to choose a filled group",
			"chose":                   " chose ",
			"group":                   "group ",

			// poll

			"table":           "Queue:\n",
			"youWereAssigned": "You were assigned the number ",
			"youAre":          ". You are ",
			"inQueue":         " positions away from the start.",
			"chooseYourGroup": "Choose your group",
			"pollText": "The poll is located below. \n\n" +
				"When you have decided, please, press the button below.\n\n",
			"button":           "button ↓",
			"choosingRightNow": "The one choosing right now is ",
			"outsideTheirTurn": " affected the poll <b>outside their choosing turn</b> ⚠️⚠️⚠️" +
				" \n\nThe poll has to be restarted now.",
			"outsideYourTurn":    "You have voted outside your turn. The poll must be restarted now.",
			"messageDuringPoll":  "Messages are not accepted during the poll.",
			"pollButton":         "Poll",
			"shutdownButton":     "Disable bot",
			"sendButton":         "Send",
			"emergencyPollStop":  "The poll was stopped.",
			"noChoice":           "You cannot choose an empty group.",
			"fullGroup":          "You cannot choose a full group.",
			"youChose":           "You have chosen ",
			"thanksForAnswering": "Thanks for answering. ",
			"successfulPollEnd": "The poll has ended. Thanks for participating. " +
				"The poll was sent to ",
			"endResults":     "The poll has ended. Results are provided above.",
			"infoText1":      "Choose a non-filled group. The counter for each group is: \n\n",
			"infoText2":      "\n\nAllowed groups are those with <b>less than ",
			"infoText3":      "</b> people. The people who have chosen a group:\n\n",
			"thanksForTime":  "Thanks for your time.",
			"waitForPollEnd": "Wait for the poll to end.",
			"badPollEnd": "Something went wrong during the poll. The poll has ended.\n\n" +
				"Register to a new one by pressing the button below.",
		},
		"russian": {
			// idle
			"idle":         "Этот бот сейчас выключен. Попросите владельца запустить бота.",
			"startBot":     "Запуск",
			"pressToStart": "Нажмите на кнопку ниже для запуска бота.",
			"startingBot":  "Окей, запускаю бота...",

			// no poll
			"register": "Регистрация",
			"quit":     "Выход",
			"start": "Добро пожаловать! \n<b>Использование:</b> \n\t- \"<b>Регистрация</b>\" для регистрации в голосовании" +
				" \n\t- \"<b>Выход</b>\", чтобы не участвовать" +
				"\n\t- \"<b>Помощь</b>\" для помощи \n\t- \"<b>Список</b>\" для получения списка участвующих.",
			"start2":                       "Выберите, будете ли участвовать в голосовании.",
			"unknownCommand":               "Неизвестная команда",
			"list":                         "Список",
			"help":                         "Помощь",
			"listCase":                     "Сейчас участвует ",
			"listCase1":                    " человека в голосовании.\n",
			"listCase2":                    " человек в голосовании.\n",
			"userEnlists":                  "Вы <b>участвуете</b> в голосовании.",
			"userAlreadyInPoll":            "Вы уже участвуете в голосовании.",
			"userAlreadyNotInPoll":         "Вы и так не участвуете в голосовании.",
			"userDeletedFromPoll":          "Вы удалены из голосования",
			"userDoesntEnlist":             "Вы <b>не</b> участвуете в голосовании.",
			"userIsThinkingAboutEnlisting": "Выберите, будете ли участвовать в голосовании.",
			"notEnoughUsers":               "Недостаточно пользователей.",
			"pollStarted":                  "Голосование началось.",
			"notPermitted":                 "Вы не имеете права на это действие.",
			"shutdown":                     "Теперь бот выключен.",
			"unknownMessage":               "Я не знаю, как ответить на это сообщение.",
			"addedToPoll":                  "Вы добавлены в голосование.",
			"someoneAddedToPoll":           " добавлен(-а) в голосование.",
			"someoneDeletedFromPoll":       " удален(-а) из голосования.",

			// helpers

			"cancelledTheirVote":      " отменил(-а) свой голос!",
			"handleBug":               "забанят",
			"groupFilled":             "Вам нужно перевыбрать. Эта группа заполнена доверху.",
			"someoneTriedFilledGroup": " попытался(-ась) выбрать заполненную группу",
			"chose":                   " выбрал(-а) ",
			"group":                   "группа ",

			// poll

			"table":           "Очередь:\n",
			"youWereAssigned": "Вам был назначен номер ",
			"youAre":          ". Вы ",
			"inQueue":         " в порядке очереди.",
			"chooseYourGroup": "Выберите группу",
			"pollText": "Опрос предоставлен ниже. \n\n" +
				"Когда сделаете окончательный выбор, пожалуйста, нажмите на кнопку ниже.\n\n",
			"button":             "кнопка ↓",
			"choosingRightNow":   "Сейчас выбирает ",
			"outsideTheirTurn":   " повлиял на голосование <b>не в свою очередь</b> ⚠️⚠️⚠️ \n\nТеперь придётся проводить голосование заново.",
			"outsideYourTurn":    "Вы проголосовали не в свою очередь. Теперь придётся провести голосование заново.",
			"messageDuringPoll":  "Сообщения в режиме голосования не принимаются.",
			"pollButton":         "Голосование",
			"shutdownButton":     "Отключить",
			"sendButton":         "Отправить",
			"emergencyPollStop":  "Голосование завершено.",
			"noChoice":           "Вы должны выбрать хоть какую-нибудь группу.",
			"fullGroup":          "Вы должны выбрать группу, в которой есть места.",
			"youChose":           "Вы выбрали ",
			"thanksForAnswering": "Спасибо за ответ. ",
			"successfulPollEnd": "Голосование завершено. Всем спасибо за участие. " +
				"Опрос отправлен ",
			"endResults":     "Голосование завершено. Результаты выше.",
			"infoText1":      "Выбирайте группу, в которой есть места. Счётчик мест для каждой группы: \n\n",
			"infoText2":      "\n\nПодходят группы, в которых <b>меньше ",
			"infoText3":      "</b> человек. Люди, выбравшие группы:\n\n",
			"thanksForTime":  "Спасибо за уделенное время.",
			"waitForPollEnd": "Дождитесь окончания голосования.",
			"badPollEnd": "Что-то пошло не так во время голосования. Голосование завершено.\n\n" +
				"Зарегистрируйтесь на новое голосование, нажав кнопку ниже.",
		},
	}
	lang = language[BotLanguage]
)
