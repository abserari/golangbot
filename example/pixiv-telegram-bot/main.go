package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/NateScarlet/pixiv/pkg/artwork"
	"github.com/NateScarlet/pixiv/pkg/client"
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

var modesetHelpString = `
try /set d
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
			cmd := strings.Split(update.CallbackQuery.Data, " ")
			if len(cmd) == 2 {
				switch cmd[0] {
				case "artwork":
					// 画作详情
					i := &artwork.Artwork{ID: cmd[1]}
					ctx := client.With(context.Background(), pixivClient.Papi)
					err := i.Fetch(ctx) // 获取画作详情(不含分页), 直接更新 struct 数据。
					if err != nil {
						logger.Error(err.Error())
					}
					if i.Image.Original == "" {
						bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "No origin image"))
						continue
					} else {
						edited := tgbotapi.EditMessageTextConfig{
							BaseEdit: tgbotapi.BaseEdit{
								InlineMessageID: update.CallbackQuery.InlineMessageID,
								ReplyMarkup:     nil,
							},
							Text: i.Image.Original,
						}
						bot.Send(edited)
					}
				}
			}

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			// reply the callback from inlinekeyboard
			log.Println(update.CallbackQuery.Message)
			// msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.InlineMessageID)
			// switch update.CallbackQuery.Data {
			// default:
			// }
			// bot.Send(msg)
			continue
		}

		if update.InlineQuery != nil {
			log.Println(update.InlineQuery.Query)
			var query string
			var nsfw = artwork.ContentRatingR18
			cmd := strings.Split(update.InlineQuery.Query, " ")
			if update.InlineQuery.Query == "" {
				query = "top"
			} else {
				query = cmd[0]
			}
			for _, v := range cmd {
				switch v {
				case "r18":
				case "nsfw":
				case "r":
				default:
					nsfw = artwork.ContentRatingSafe
				}
			}
			if len(cmd) == 2 {
				query = cmd[1]
			}

			if update.InlineQuery.Offset == "" {
				update.InlineQuery.Offset = "1"
			}

			offset, err := strconv.Atoi(update.InlineQuery.Offset)
			if err != nil {
				logger.Error("not int offset")
			}
			photos, err := pixivClient.Search(context.Background(), query, artwork.OrderDateASC, nsfw, artwork.SearchModePartialTag, offset)
			if err != nil {
				logger.Error(err.Error())
			}

			resultPhotos := make([]interface{}, 0)
			// for i, v := range photos {
			// 	if i > 4 {
			// 		break
			// 	}
			// 	picture := tgbotapi.NewInlineQueryResultPhotoWithThumb(strconv.Itoa(i+offset)+update.InlineQuery.Query, v.Image.Thumb, v.Image.Thumb)
			// 	picture.MimeType = "image/jpeg"
			// 	picture.Description = v.Description
			// 	picture.Title = v.Title
			// 	captionFormat := "<a href='%s'>%s</a>\nUser: <a href='%s'>%s</a>"
			// 	picture.Caption = fmt.Sprintf(captionFormat, "https://www.pixiv.net/artworks/"+v.ID, v.Title, "https://www.pixiv.net/en/users/"+v.Author.ID, v.Author.Name)
			// 	picture.ParseMode = "HTML"

			// 	resultPhotos = append(resultPhotos, picture)
			// }

			for i, v := range photos {
				if i > 4 {
					break
				}
				picture := tgbotapi.NewInlineQueryResultArticle(strconv.Itoa(i+offset)+update.InlineQuery.Query, "hello", "hello")
				// picture.MimeType = "image/jpeg"
				picture.ThumbURL = v.Image.Thumb
				picture.Description = v.Description
				picture.Title = v.Title
				captionFormat := "<a href='%s'>%s</a>\nUser: <a href='%s'>%s</a>"

				s := tgbotapi.InputTextMessageContent{}
				s.Text = fmt.Sprintf(captionFormat, "https://www.pixiv.net/artworks/"+v.ID, v.Title, "https://www.pixiv.net/en/users/"+v.Author.ID, v.Author.Name)
				s.ParseMode = "HTML"
				picture.InputMessageContent = s

				picture.ReplyMarkup = inlineQueryInlineKeyboard(v.ID)

				resultPhotos = append(resultPhotos, picture)
			}

			inlineConf := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       resultPhotos,
				NextOffset:    strconv.Itoa(len(resultPhotos) - 1 + offset),
			}

			apiResp, err := bot.AnswerInlineQuery(inlineConf)
			if err != nil {
				log.Println(err)
			}
			log.Println(apiResp)
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		// Command Handle
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "help":
				msg.Text = `try /open or /close 
use /set to see which mode could be set
use /top to see top 10 pictures`
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
					msg.Text = fmt.Sprintf("now mode %s", mode)
				default:
					msg.Text = CurrentState(update.Message.From.ID)
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

func inlineQueryInlineKeyboard(id string) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Origin", "artwork "+id),
		),
	)
	return &keyboard
}

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

func CurrentState(id int) string {
	return strings.Replace(modesetHelpString,
		pages[id].mode,
		fmt.Sprintf("%s (default)", pages[id].mode), 1)
}
