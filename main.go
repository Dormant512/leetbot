package main

import (
	"context"
	"fmt"
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

var statsMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("View on Leetcode", "go2user")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "main")),
)

var randomMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Easy", "easy")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Medium", "medium")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Hard", "hard")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "main")),
)

var easyMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("View on Leetcode", "go2easy")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "random")),
)

var mediumMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("View on Leetcode", "go2medium")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "random")),
)

var hardMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("View on Leetcode", "go2hard")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "random")),
)

var dailyMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("View on Leetcode", "go2daily")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "main")),
)

var aboutMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("View source", "go2source")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "main")),
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
	var easyURL, mediumURL, hardURL, dailyURL, userURL string
	sourceURL := "https://github.com/Dormant512/leetbot"
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
				if err == nil && respData.ActiveDailyCodingChallengeQuestion.Date != "" {
					task := respData.ActiveDailyCodingChallengeQuestion
					q := task.Question
					dailyText = "*Daily task for " + task.Date + "*\n\n" + q.Title + "\n\n*Difficulty*: " + q.Difficulty
					dailyText += "\n" + "*Acceptance rate*: " + strconv.FormatFloat(q.AcRate, 'f', 0, 64)
					dailyText += "%\n" + "*Tags*:"
					for _, val := range q.TopicTags {
						dailyText += "\n‣ " + val.Name
					}
				}
				postMessage(update, dailyText, dailyMenu, bot)

			case "random":
				// Handle random task
				postMessage(update, "Pick difficulty:", randomMenu, bot)

			case "about":
				// Handle about the bot
				postMessage(update, aboutText, aboutMenu, bot)

			// SECOND LEVEL MENU
			case "easy":
				// Handle easy task
				easyMessage := HandleTask("easy")
				postMessage(update, easyMessage, easyMenu, bot)

			case "medium":
				// Handle medium task
				mediumMessage := HandleTask("medium")
				postMessage(update, mediumMessage, mediumMenu, bot)

			case "hard":
				// Handle hard task
				hardMessage := HandleTask("hard")
				postMessage(update, hardMessage, hardMenu, bot)

			case "go2user":
				// Handle user URL
				fmt.Println("Went to user " + userURL)

			case "go2easy":
				// Handle easy URL
				fmt.Println("Went to user " + easyURL)

			case "go2medium":
				// Handle medium URL
				fmt.Println("Went to user " + mediumURL)

			case "go2hard":
				// Handle hard URL
				fmt.Println("Went to user " + hardURL)

			case "go2daily":
				// Handle daily URL
				fmt.Println("Went to user " + dailyURL)

			case "go2source":
				// Handle source URL
				fmt.Println("Went to user " + sourceURL)

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
				postMessage(update, statText, statsMenu, bot)
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
