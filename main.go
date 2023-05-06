package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/machinebox/graphql"
	"log"
	"os"
	"strconv"
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

var greet = `Hello, I am a utility LeetCode bot.
You can:
- view your LeetCode stats,
- solve the daily task,
- pick a random one,
- read about the bot.
Your choice?`

var aboutText = `This is a simple Telegram bot implemented in Go for an internship in VK.
It levers the Telegram Bot API:
github.com/go-telegram-bot-api/telegram-bot-api/v5
Are you interested in the source code?`

//var requestDaily string = `query questionOfToday {
//   activeDailyCodingChallengeQuestion {
//       date
//       userStatus
//       link
//       question {
//           acRate
//           difficulty
//           freqBar
//           frontendQuestionId: questionFrontendId
//           isFavor
//           paidOnly: isPaidOnly
//           status
//           title
//           titleSlug
//           hasVideoSolution
//           hasSolution
//           topicTags {
//               name
//               id
//               slug
//           }
//       }
//   }
//}`

var requestUser string = `query getUserProfile($username: String!) {
	matchedUser(username: $username) {
		username
		submitStats: submitStatsGlobal {
			acSubmissionNum {
				difficulty
				count
				submissions
			}
		}
	}
}`

type UserData struct {
	MatchedUser UserStats `json:"matchedUser"`
}

type UserStats struct {
	Username    string `json:"username"`
	SubmitStats struct {
		AcSubmissionNum []struct {
			Difficulty  string `json:"difficulty"`
			Count       int    `json:"count"`
			Submissions int    `json:"submissions"`
		} `json:"acSubmissionNum"`
	} `json:"submitStats"`
}

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
	var promptName bool

	for update := range updates {
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery.Data

			switch callback {

			// FIRST LEVEL MENU
			case "stats":
				// Handle user stats
				promptName = true
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Enter username:")
				bot.Send(msg)

			case "daily":
				// Handle daily task
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "DAILY HERE.")
				msg.ReplyMarkup = dailyMenu
				bot.Send(msg)

			case "random":
				// Handle random task
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Pick difficulty:")
				msg.ReplyMarkup = randomMenu
				bot.Send(msg)

			case "about":
				// Handle about the bot
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, aboutText)
				msg.ReplyMarkup = aboutMenu
				bot.Send(msg)

			// SECOND LEVEL MENU
			case "easy":
				// Handle easy task
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "EASY HERE.")
				msg.ReplyMarkup = easyMenu
				bot.Send(msg)

			case "medium":
				// Handle medium task
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "MEDIUM HERE.")
				msg.ReplyMarkup = mediumMenu
				bot.Send(msg)

			case "hard":
				// Handle hard task
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "HARD HERE.")
				msg.ReplyMarkup = hardMenu
				bot.Send(msg)

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
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Your choice?")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)

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
					statText = "Stats for user " + respData.MatchedUser.Username
					for _, val := range respData.MatchedUser.SubmitStats.AcSubmissionNum {
						statText += "\n" + val.Difficulty + ": " + strconv.Itoa(val.Count) + " tasks"
					}
				}

				promptName = false
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, statText)
				msg.ReplyMarkup = statsMenu
				bot.Send(msg)
				continue
			}
			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, greet)
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your choice?")
			msg.ReplyMarkup = mainMenu
			bot.Send(msg)
		}
	}
}
