package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yhyddr/golangbot/eval"
)

var maincode = `
package main

%s
`
var fmtcode = `
package main

import (
	"fmt"
)

func main() {
	%s
}
`
var setcodemode = "```go %s```"

func init() {
	os.Setenv("TELEGRAM_TECHCATS_BOT_TOKEN", "THIS IS YOUR TEMP ID")
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TECHCATS_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// update is every new message
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Println(update.Message.Chat.ID, "come on")

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = `type /eval fmt.Println("Hello, World") : now only for fmt package.
type /run import ("fmt")
				func main() {
					fmt.Println("Hello, World")
				}
				 : using like go playground, but don't need package main `
			case "run":
				{
					code := strings.NewReplacer(`“`, `"`, `”`, `"`).Replace(update.Message.CommandArguments())
					res, err := eval.GoCode(fmt.Sprintf(maincode, code))
					if err != nil {
						log.Println(err)
						continue
					}
					if res.Errors != "" {
						msg.Text = res.Errors
					} else {
						for _, e := range res.Events {
							if e.Kind == "stdout" {
								msg.Text = fmt.Sprintf(setcodemode, e.Message)
								msg.ParseMode = "MarkdownV2"
								continue
							}
						}
					}
				}
			case "eval":
				{
					// handle code
					code := strings.NewReplacer(`“`, `"`, `”`, `"`).Replace(update.Message.CommandArguments())
					res, err := eval.GoCode(fmt.Sprintf(fmtcode, code))
					if err != nil {
						log.Println(err)
						continue
					}
					if res.Errors != "" {
						msg.Text = res.Errors
					} else {
						for _, e := range res.Events {
							if e.Kind == "stdout" {
								msg.Text = e.Message
								continue
							}
						}
					}
				}
			default:
				msg.Text = "I don't know that command Try /help"
			}
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
			continue
		}

		// just repeat msg
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// reply would refer the message
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
