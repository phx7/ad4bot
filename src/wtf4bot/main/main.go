package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/subosito/gotenv"
	"log"
	"os"
)

func main() {
	gotenv.Load()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Sticker != nil {
			log.Printf("Deleted sticker from @%s (%s %s)", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName)
			bot.DeleteMessage(tgbotapi.DeleteMessageConfig{update.Message.Chat.ID, update.Message.MessageID})
		}
	}
}
