package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
	"github.com/vence722/convert"
	"gopkg.in/telegram-bot-api.v4"
)

var replyChatIDs = []int64{}

func StartTelegramBot(crawler func() ([]*Match, error)) error {
	bot, err := tgbotapi.NewBotAPI("627688442:AAHjNsFHqzc522NADbBgAzxRdBGdWg1hZ4g")
	if err != nil {
		return err
	}

	startCronJob(bot, crawler)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		inMsg := update.Message.Text
		fmt.Println("incoming message:", inMsg)

		if inMsg == "/start" || inMsg == "/help" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to use Vence football news bot!\ncommands:\n/help show help message\n/latest show latest updates\n/subscribe subscribe updates\n/unsubscribe unsubscribe updates")
			bot.Send(msg)
		} else if inMsg == "/latest" {
			matches, err := crawler()
			if err != nil {
				fmt.Println("Something wrong when crawling match results")
			}
			if len(matches) == 0 {
				break
			}
			// format and send matches results
			content := formatMatches(matches)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
			bot.Send(msg)
		} else if inMsg == "/subscribe" {
			// add chat id in reply list
			replyChatIDs = append(replyChatIDs, update.Message.Chat.ID)
			fmt.Println("User", update.Message.From.UserName, "subscribed")

			// reply message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Thanks for subscribing this channel. Updates for football matches results will be sent to you at 9:00 AM every day.")
			bot.Send(msg)
		} else if inMsg == "/unsubscribe" {
			// remove chat id in replay list
			for i, id := range replyChatIDs {
				if id == update.Message.Chat.ID {
					replyChatIDs = append(replyChatIDs[:i], replyChatIDs[i+1:]...)
					break
				}
			}
			fmt.Println("User", update.Message.From.UserName, "unsubscribed")

			// reply message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You have unsubscribed from this channel")
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I don't know what you're saying.")
			bot.Send(msg)
		}
	}

	return nil
}

func startCronJob(bot *tgbotapi.BotAPI, crawler func() ([]*Match, error)) {
	c := cron.New()
	c.AddFunc("0 0 9 * * *", func() {
		matches, err := crawler()
		if err != nil {
			fmt.Println("Something wrong when crawling match results")
		}
		if len(matches) == 0 {
			return
		}
		for _, replyChatID := range replyChatIDs {
			// format and send matches results
			content := formatMatches(matches)
			msg := tgbotapi.NewMessage(replyChatID, content)
			bot.Send(msg)
		}
	})
	c.Start()
}

func formatMatches(matches []*Match) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// format matches
	content := "Matches results at " + time.Now().In(loc).Format("2006-01-02 15:04:05") + ":\n"
	for _, match := range matches {
		content += (match.Competition + "|" + match.Time.In(loc).Format("2006-01-02 15:04:05") + "|" + match.HomeTeam + " " + convert.Int2Str(match.HomeScore) + ":" + convert.Int2Str(match.AwayScore) + " " + match.AwayTeam + "\n")
	}
	return content
}
