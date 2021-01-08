package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/abserari/golangbot/pkg/pixiv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// func init() {
// 	os.Setenv("TELEGRAM_TECHCATS_BOT_TOKEN", "THIS IS YOUR TEMP ID")
// }

var inlineTopKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("进入班城剧组", "https://www.yuque.com/bandle"),
		tgbotapi.NewInlineKeyboardButtonSwitch("转发 /remake", "/remake"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("上一页", "pageleft"),
		tgbotapi.NewInlineKeyboardButtonData("下一页", "pageright"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("上一天", "dateup"),
		tgbotapi.NewInlineKeyboardButtonData("下一天", "datedown"),
	),
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButtonContact("Contact"),
		tgbotapi.NewKeyboardButtonLocation("location"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
	}
	// Reading app id from env (never hardcode it!).
	botToken := viper.GetString("TgBotToken")

	cookies := viper.GetString("pixiv.Cookies")

	logger, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.InfoLevel))
	defer func() { _ = logger.Sync() }()

	config := &pixiv.ClientConfig{
		Cookies: cookies,
		Logger:  logger,
	}

	pixivClient := pixiv.NewClient(config)
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	// update is every new message
	for update := range updates {
		if update.CallbackQuery != nil {
			// reply the callback from inlinekeyboard
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			switch update.CallbackQuery.Data {
			default:
			}
			bot.Send(msg)
			continue
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// Command Handle
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "help":
				msg.Text = `try /top`
			case "top":
				{
					go func(update tgbotapi.Update) {
						cmd := strings.Split(update.Message.Text, " ")
						log.Println(cmd)
						page := 0
						mode := ""
						switch len(cmd) {
						case 3:
							i, err := strconv.Atoi(cmd[2])
							if err != nil {
								// not pages
							}
							page = i

							if cmd[1] == "nsfw" {
								mode = "daily_r18"
							} else {
								mode = cmd[1]
							}
						case 2:
							i, err := strconv.Atoi(cmd[1])
							if err != nil {
								// not pages
							}
							page = i
						default:

						}
						items, err := pixivClient.Ranking(context.Background(), mode, "", "2020-01-04", page)
						if err != nil {
							logger.Info(err.Error())
							msg.Text = err.Error()
							if _, err := bot.Send(msg); err != nil {
								log.Println(err)
							}
							return
						}
						var mediagroup = make([]interface{}, 0)
						for i := 0; i < len(items); i++ {
							mediagroup = append(mediagroup,
								tgbotapi.NewInputMediaPhoto(items[i].Image.Regular))
							if i >= 9 {
								break
							}
						}

						cfg := tgbotapi.NewMediaGroup(msg.ChatID, mediagroup)
						cfg.ReplyMarkup = inlineTopKeyboard
						if _, err := bot.Send(cfg); err != nil {
							log.Println(err)
						}
					}(update)
					continue
				}
			default:
				msg.Text = "I don't know that command Try /help"
			}
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
			continue
		}
	}
}
