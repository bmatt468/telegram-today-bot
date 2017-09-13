package todaybot

import (
	"github.com/mmcdole/gofeed"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
	"time"
)

type TodayBot struct {
	API *tgbotapi.BotAPI
}

func Connect(tkn string, debug bool) (*TodayBot, error) {
	bot, err := tgbotapi.NewBotAPI(tkn)

	if err != nil {
		log.Fatal(err)
	}

	if debug {
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}

	bot.Debug = debug
	tbot := &TodayBot{
		API: bot,
	}

	return tbot, err
}

func (bot *TodayBot) OpenWebhook(url string) {
	_, err := bot.API.SetWebhook(tgbotapi.NewWebhook(url))
	if err != nil {
		log.Fatal(err)
	}
}

func (bot *TodayBot) Listen(token string) <-chan tgbotapi.Update {
	updates := bot.API.ListenForWebhook(token)
	return updates
}

func (bot *TodayBot) ParseAndExecuteUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		switch cmd := strings.Split(update.Message.Text, " "); strings.Replace(cmd[0], "@TodaysHolidaysBot", "", -1) {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Good Morning!")
			bot.API.Send(msg)

		case "/today":
			text := "Today is " + time.Now().Format("Monday, January 2, 2006\n")
			text += "Today's Holidays:\n"

			fp := gofeed.NewParser()
			holidays, _ := fp.ParseURL("https://www.checkiday.com/rss.php?tz=America/New_York")
			for _, holiday := range holidays.Items {
				text += holiday.Title + "\n"
				//text += fmt.Sprintf("[%s](%s)\n",holiday.Title, holiday.Link)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.API.Send(msg)
		}
	}
}
