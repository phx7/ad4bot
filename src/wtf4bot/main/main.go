package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func main() {
	// init
	godotenv.Load(".env")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account @%s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// delete stickers
		if (update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup()) && update.Message.Sticker != nil {
			log.Printf("Deleted sticker from @%s (%s %s)", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName)
			bot.DeleteMessage(tgbotapi.DeleteMessageConfig{update.Message.Chat.ID, update.Message.MessageID})
		}

		// show rules to newcomers
		if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			if update.Message.NewChatMembers != nil {
				var newUsers []string

				for _, user := range *update.Message.NewChatMembers {
					newUsers = append(newUsers, "@"+getUserName(user))
				}

				joinedUsers := strings.Join(newUsers, " ")

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Приветствую, %s\nВ нашем чате необходимо соблюдать правила, посмотреть их ты можешь тут %s", joinedUsers, os.Getenv("RULES_LINK")))
				send(bot, msg)
			}
		}

		// COMMANDS
		if update.Message.IsCommand() && (update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup()) {

			switch update.Message.Command() {
			case "rules":
				user, _ := bot.GetChatMember(tgbotapi.ChatConfigWithUser{update.Message.Chat.ID, update.Message.Chat.UserName, update.Message.From.ID})
				if user.IsCreator() {
					rules_message_link := update.Message.CommandArguments()
					os.Setenv("RULES_LINK", rules_message_link)
					writeConfig()
					resp := tgbotapi.NewMessage(update.Message.Chat.ID, "Ссылка на сообщение с правилами установлена: "+rules_message_link)
					resp.ParseMode = tgbotapi.ModeHTML
					send(bot, resp)
				}
			}
		}
	}
}

func getUserName(user tgbotapi.User) string {
	if user.UserName == "" {
		return user.FirstName
	}
	return user.UserName
}

func send(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("[FAILED_TO_SEND_MESSAGE] [%v]", msg)
	}
}

func writeConfig() {
	env := make(map[string]string)
	env["BOT_TOKEN"] = os.Getenv("BOT_TOKEN")
	env["RULES_LINK"] = os.Getenv("RULES_LINK")
	env["HTTP_PROXY"] = os.Getenv("HTTP_PROXY")
	godotenv.Write(env, ".env")
	godotenv.Load(".env")
}
