package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

// func init() {
// 	os.Setenv("TELEGRAM_TECHCATS_BOT_TOKEN", "THIS IS YOUR TEMP ID")
// }

var inlineNumericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("进入班城剧组", "https://www.yuque.com/bandle"),
		tgbotapi.NewInlineKeyboardButtonSwitch("转发 /remake", "/remake"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("yes", "/yes"),
		tgbotapi.NewInlineKeyboardButtonData("no", "/no"),
	),
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
		// handle error
	}

	return string(body), err
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TECHCATS_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	var PeopleCount int
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// update is every new message
	for update := range updates {
		if update.CallbackQuery != nil {
			// just repeat the callback
			fmt.Print(update)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			switch update.CallbackQuery.Data {

			case "/yes":
				{
					if PeopleCount > 4 {
						msg.Text = `Remaking`
						PeopleCount = 0
						break
					}
					msg.Text = `您是否要发起投降, 当前投降人数 ` + strconv.Itoa(PeopleCount) + `/5. \n 输入 /remake 或 /yes 同意, /no 拒绝`

				}
			case "/no":
				{
					msg.Text = `您已经拒绝投降`

				}
			}
			bot.Send(msg)
			continue
		}
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Println("come on message", update.CallbackQuery)

		if update.Message.Text == "/诗歌" {
			go func(update tgbotapi.Update) {
				var err error
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text, err = httpGet("https://v1.jinrishici.com/rensheng.txt")
				if err != nil {
					fmt.Println(err)
				}

				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
				}
			}(update)

			continue
		}

		// Command Handle
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "nsfw":
				{
					go func(update tgbotapi.Update) {
						var err error
						msg.Text, err = httpGet("http://rakuen.thec.me/PixivRss/male_r18-20")
						if err != nil {
							fmt.Println(err)
						}

						if _, err := bot.Send(msg); err != nil {
							log.Println(err)
						}
					}(update)

					continue
				}
			case "help":
				msg.Text = `type /eval fmt.Println("Hello, World") : now only for fmt package.
type /run import ("fmt")
				func main() {
					fmt.Println("Hello, World")
				}
				 : using like go playground, but don't need package main 
				 
				 If you use Phone or Compatible APP, you could type /open and /close to open and close a keyboard.
				 `
			case "remake":
				{
					PeopleCount++
					if PeopleCount > 4 {
						msg.Text = `Remaking`
						PeopleCount = 0
						break
					}
					msg.Text = `您是否要发起投降, 当前投降人数 ` + strconv.Itoa(PeopleCount) + `/5. \n 输入 /remake 或 /yes 同意, /no 拒绝`
					msg.ReplyMarkup = inlineNumericKeyboard
				}
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
			case "open":
				msg.ReplyMarkup = numericKeyboard
			case "close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			default:
				msg.Text = "I don't know that command Try /help"
			}
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
			continue
		}

		// // just repeat msg
		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// // reply would refer the message
		// msg.ReplyToMessageID = update.Message.MessageID

		// if _, err := bot.Send(msg); err != nil {
		// 	log.Println(err)
		// }
	}
}
