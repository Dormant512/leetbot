package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/machinebox/graphql"
	"log"
	"os"
	"strconv"
	"time"
)

var mainMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("User stats", "stats")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Daily task", "daily")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Random task", "random")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("About", "about")),
)

var randomMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Easy", "easy")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Medium", "medium")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Hard", "hard")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "main")),
)

var urlMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("GOTO URL", "go2url")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "random")),
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	promptName := false

	for update := range updates {
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery.Data

			switch callback {

			// FIRST LEVEL MENU
			case "stats":
				// Handle user stats
				promptName = true
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Enter username:")
				msg.ParseMode = tgbotapi.ModeMarkdown
				bot.Send(msg)

			case "daily":
				// Handle daily task
				var respData DailyData
				ctx := context.Background()
				client := graphql.NewClient("https://leetcode.com/graphql")
				req := graphql.NewRequest(requestDaily)
				err := client.Run(ctx, req, &respData)

				dailyText := "Daily task for " + time.Now().Format("02.01.2006") + " *not* found."
				if err != nil || respData.ActiveDailyCodingChallengeQuestion.Date == "" {
					continue
				}
				task := respData.ActiveDailyCodingChallengeQuestion
				q := task.Question
				dailyText = "*Daily task for " + task.Date + "*\n\n" + q.Title + "\n\n*Difficulty*: " + q.Difficulty
				dailyText += "\n" + "*Acceptance rate*: " + strconv.FormatFloat(q.AcRate, 'f', 0, 64)
				dailyText += "%\n" + "*Tags*:"
				for _, val := range q.TopicTags {
					dailyText += "\n‣ " + val.Name
				}

				url := "https://leetcode.com/problems/" + q.TitleSlug + "/"
				cbd := "main"
				thisMenu := urlMenu
				thisMenu.InlineKeyboard[0][0].Text = "View on Leetcode"
				thisMenu.InlineKeyboard[0][0].URL = &url
				thisMenu.InlineKeyboard[1][0].CallbackData = &cbd
				postMessage(update, dailyText, thisMenu, bot)

			case "random":
				// Handle random task
				postMessage(update, "Pick difficulty:", randomMenu, bot)

			case "about":
				// Handle about the bot
				url := "https://github.com/Dormant512/leetbot"
				cbd := "main"
				thisMenu := urlMenu
				thisMenu.InlineKeyboard[0][0].Text = "View source on GitHub"
				thisMenu.InlineKeyboard[0][0].URL = &url
				thisMenu.InlineKeyboard[1][0].CallbackData = &cbd
				postMessage(update, aboutText, thisMenu, bot)

			// SECOND LEVEL MENU
			case "easy":
				// Handle easy task
				mes, url := HandleTask("easy")
				thisMenu := urlMenu
				thisMenu.InlineKeyboard[0][0].Text = "View on Leetcode"
				thisMenu.InlineKeyboard[0][0].URL = &url
				postMessage(update, mes, thisMenu, bot)

			case "medium":
				// Handle medium task
				mes, url := HandleTask("medium")
				thisMenu := urlMenu
				thisMenu.InlineKeyboard[0][0].Text = "View on Leetcode"
				thisMenu.InlineKeyboard[0][0].URL = &url
				postMessage(update, mes, thisMenu, bot)

			case "hard":
				// Handle hard task
				mes, url := HandleTask("hard")
				thisMenu := urlMenu
				thisMenu.InlineKeyboard[0][0].Text = "View on Leetcode"
				thisMenu.InlineKeyboard[0][0].URL = &url
				postMessage(update, mes, thisMenu, bot)

			case "main":
				// Handle back to main
				postMessage(update, "Your choice?", mainMenu, bot)

			default:
				log.Printf("Unknown callback query: %s", callback)
			}
		} else if update.Message != nil {
			if promptName {
				username := update.Message.Text

				var respData UserData
				ctx := context.Background()
				client := graphql.NewClient("https://leetcode.com/graphql")
				req := graphql.NewRequest(requestUser)
				req.Var("username", username)
				err := client.Run(ctx, req, &respData)

				statText := "User " + username + " not found."
				if err == nil && respData.MatchedUser.Username != "" {
					statText = "*Stats for user " + respData.MatchedUser.Username + "*\n"
					for _, val := range respData.MatchedUser.SubmitStats.AcSubmissionNum {
						statText += "\n*" + val.Difficulty + "*: " + strconv.Itoa(val.Count) + " tasks"
					}
				}

				promptName = false
				url := "https://leetcode.com/" + username + "/"
				cbd := "main"
				thisMenu := urlMenu
				thisMenu.InlineKeyboard[0][0].Text = "View on Leetcode"
				thisMenu.InlineKeyboard[0][0].URL = &url
				thisMenu.InlineKeyboard[1][0].CallbackData = &cbd

				if err != nil || respData.MatchedUser.Username == "" {
					thisMenu.InlineKeyboard = thisMenu.InlineKeyboard[1:]
				}

				postMessage(update, statText, thisMenu, bot)
				continue
			}
			if update.Message.Text == "/start" {
				postMessage(update, greet, mainMenu, bot)
				continue
			}

			postMessage(update, "Your choice?", mainMenu, bot)
		}
	}
}

func postMessage(update tgbotapi.Update, message string, menu tgbotapi.InlineKeyboardMarkup, bot *tgbotapi.BotAPI) {
	var msg tgbotapi.MessageConfig
	if update.CallbackQuery != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
	} else if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, message)
	} else {
		return
	}
	msg.ReplyMarkup = menu
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)
}
