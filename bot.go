package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yhyddr/golangbot/eval"
)

var fmtcode = `
package main

import (
	"fmt"
)

func main() {
	%s
}
`

func init() {
	os.Setenv("TELEGRAM_APITOKEN", "969298533:AAF_AE2DJfjlcVUxy1U44B185rKIFrfwzUM")
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
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

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = `type /eval fmt.Println("Hello, World") now only for fmt package.
				type /run package main

				import (
					"fmt"
				)
				
				func main() {
					fmt.Println("Hello, World")
				}. use like go playground`
			case "run":
				{
					code := strings.NewReplacer(`“`, `"`, `”`, `"`).Replace(update.Message.CommandArguments())
					res, err := eval.GoCode(code)
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
				msg.Text = "I don't know that command"
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
