package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/abserari/golangbot/pkg/pixiv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type userState struct {
	pages int
	time  time.Time
	mode  string
}

var pixivClient *pixiv.Client
var updates tgbotapi.UpdatesChannel
var bot *tgbotapi.BotAPI
var logger, _ = zap.NewDevelopment(zap.IncreaseLevel(zapcore.InfoLevel))
var pages = make(map[int]userState)

func init() {
	defer func() { _ = logger.Sync() }()

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
	log.Println(botToken)
	cookies := viper.GetString("pixiv.Cookies")

	config := &pixiv.ClientConfig{
		Cookies: cookies,
		Logger:  logger,
	}

	pixivClient = pixiv.NewClient(config)
	var err error
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ = bot.GetUpdatesChan(u)
}

func main() {
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

		if update.InlineQuery != nil {
			log.Println(update.InlineQuery)
			cmd := strings.Split(update.Message.Text, " ")
			var nsfw = "r18"
			var query = cmd[0]
			for _, v := range cmd {
				switch v {
				case "r18":
				case "nsfw":
				case "r":
				default:
					nsfw = "safe"
				}
			}
			if len(cmd) == 2 {
				query = cmd[1]
			}

			if update.InlineQuery.Offset == "" {
				update.InlineQuery.Offset = "1"
			}

			i, err := strconv.Atoi(update.InlineQuery.Offset)
			if err != nil {
				logger.Error("not int offset")
			}
			photos , err := pixivClient.Search(context.Background(), query, "date_d", nsfw, "s_tag", i)
			if err != nil {

			}

			for i, v := range photos {
				thumb = InputWebDocument(img['thumb_url'], 0, 'image/jpeg', [])
				content = InputWebDocument(img['url'], 0, 'image/jpeg', [])
				"<a href='{img['url']}'>{img['title']}</a>\nUser: <a href='{img['user_link']}'>{img['user_name']}</a>"
			}
			tgbotapi.Upload
			pictures := tgbotapi.NewInlineQueryResultPhotoWithThumb(update.InlineQuery.ID, )
			pictures.Description = update.InlineQuery.Query

			inlineConf := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       []interface{}{pictures},
			}

			if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
				log.Println(err)
			}
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		// Command Handle
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "help":
				msg.Text = `try /open or /close \n use /helpmode to see which mode could be set`
			case "helpset":
				msg.Text = `
				try /mode d
				// 	d - daily (default)
				//  w - weekly
				//  m - monthly
				//  r - rookie
				//  o - original
				//  ma - male
				//  fe - female
				//  dr - daily_r18
				//  wr - weekly_r18
				//  mr - male_r18
				//  fr - female_r18
				//  rg - r18g`
			case "set":
				cmd := strings.Split(update.Message.Text, " ")
				switch len(cmd) {
				case 2:
					mode := cmd[1]
					switch mode {
					case "d":
						mode = "daily"
					case "w":
						mode = "weekly"
					case "m":
						mode = "monthly"
					case "r":
						mode = "rookie"
					case "o":
						mode = "original"
					case "ma":
						mode = "male"
					case "fe":
						mode = "female"
					case "dr":
						mode = "daily_r18"
					case "wr":
						mode = "weekly_r18"
					case "mr":
						mode = "male_r18"
					case "fr":
						mode = "female_r18"
					case "rg":
						mode = "r18g"
					default:
						mode = "daily"
					}
					state := pages[update.Message.From.ID]
					state.mode = mode
					pages[update.Message.From.ID] = state
					msg.Text = mode
				case 1:
					msg.Text = "No mode set try /helpmode"
				default:
					msg.Text = "wrong mode set try /helpmode"
				}
			case "top":
				{
					go SendRanking(update, msg, pages[update.Message.From.ID].time, pages[update.Message.From.ID].pages, pages[update.Message.From.ID].mode)
					continue
				}
			case "open":
				msg.Text = "open board"
				msg.ReplyMarkup = numericKeyboard
			case "close":
				msg.Text = "close board"
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			case "pageleft":
				if pages[update.Message.From.ID].pages <= 1 {
					state := pages[update.Message.From.ID]
					state.pages = 1
					if state.time.IsZero() {
						state.time = time.Now().Local().AddDate(0, 0, -1)
					}
					pages[update.Message.From.ID] = state
				}
				state := pages[update.Message.From.ID]
				state.pages--
				pages[update.Message.From.ID] = state
				go SendRanking(update, msg, pages[update.Message.From.ID].time, pages[update.Message.From.ID].pages, pages[update.Message.From.ID].mode)
				continue
			case "pageright":
				if pages[update.Message.From.ID].pages <= 0 {
					state := pages[update.Message.From.ID]
					state.pages = 1
					if state.time.IsZero() {
						state.time = time.Now().Local().AddDate(0, 0, -1)
					}
					pages[update.Message.From.ID] = state
				}
				state := pages[update.Message.From.ID]
				state.pages++
				pages[update.Message.From.ID] = state
				go SendRanking(update, msg, pages[update.Message.From.ID].time, pages[update.Message.From.ID].pages, pages[update.Message.From.ID].mode)
				continue
			case "dateup":
				state := pages[update.Message.From.ID]
				state.pages = 1
				if state.time.IsZero() {
					state.time = time.Now().Local().AddDate(0, 0, -1)
				}
				state.time = state.time.Local().AddDate(0, 0, 1)
				pages[update.Message.From.ID] = state
				go SendRanking(update, msg, pages[update.Message.From.ID].time, pages[update.Message.From.ID].pages, pages[update.Message.From.ID].mode)
				continue
			case "datedown":
				state := pages[update.Message.From.ID]
				state.pages = 1
				if state.time.IsZero() {
					state.time = time.Now().Local().AddDate(0, 0, -1)
				}
				state.time = state.time.Local().AddDate(0, 0, -1)
				pages[update.Message.From.ID] = state
				go SendRanking(update, msg, pages[update.Message.From.ID].time, pages[update.Message.From.ID].pages, pages[update.Message.From.ID].mode)
				continue
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

func SendRanking(update tgbotapi.Update, msg tgbotapi.MessageConfig, date time.Time, page int, mode string) {
	items, err := pixivClient.Ranking(context.Background(), mode, "", date, page)
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
	if _, err := bot.Send(cfg); err != nil {
		log.Println(err)
	}
}

// var inlineTopKeyboard = tgbotapi.NewInlineKeyboardMarkup(
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonURL("进入班城剧组", "https://www.yuque.com/bandle"),
// 		tgbotapi.NewInlineKeyboardButtonSwitch("转发今日图片", "/top"),
// 	),
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonData("上一页", "pageleft"),
// 		tgbotapi.NewInlineKeyboardButtonData("下一页", "pageright"),
// 	),
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonData("上一天", "dateup"),
// 		tgbotapi.NewInlineKeyboardButtonData("下一天", "datedown"),
// 	),
// )

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/pageleft"),
		tgbotapi.NewKeyboardButton("/pageright"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/dateup"),
		tgbotapi.NewKeyboardButton("/datedown"),
	),
)
