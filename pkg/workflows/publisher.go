package workflows

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/abserari/telegraph"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
)

func SendMessage(ctx context.Context, message interface{}) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("message", message)
	return nil
}

var account = telegraph.Account{
	ShortName:   viper.GetString("telegraph.ShortName"),   // required
	AccessToken: viper.GetString("telegraph.AccessToken"), // required

	// Author name/link can be epmty. So secure. Much anonymously. Wow.
	AuthorName: viper.GetString("telegraph.AuthorName"), // optional
	AuthorURL:  viper.GetString("telegraph.AuthorURL"),  // optional
}

var sockAddr = viper.GetString("Socks5.Address")
var socks5Pwd = viper.GetString("Socks5.Password")
var socks5Username = viper.GetString("Socks5.Username")

var botToken = viper.GetString("TgBotToken")
var chatID = viper.GetInt64("TechCatsPubChatID")

type TgBot struct {
	*tgbotapi.BotAPI
	chatID int64
}

var once sync.Once

var Tg *TgBot

func GetTgBot() *TgBot {
	once.Do(func() {
		Tg = NewTgBot()
	})
	return Tg
}

func NewTgBot() *TgBot {
	sock5, err := proxy.SOCKS5("tcp", sockAddr, &proxy.Auth{
		User:     socks5Username,
		Password: socks5Pwd,
	}, proxy.Direct)

	dc := sock5.(interface {
		DialContext(ctx context.Context, network, addr string) (net.Conn, error)
	})
	httpClient := http.DefaultClient
	if sockAddr != "" {
		httpClient.Transport = &http.Transport{
			DialContext: dc.DialContext,
		}
	}

	log.Println(sockAddr)
	client, err := tgbotapi.NewBotAPIWithClient(botToken, httpClient)
	if err != nil {
		log.Fatalln(err)
	}

	client.Debug = true

	return &TgBot{
		client,
		chatID,
	}
}

func (tg *TgBot) SendMessageToTelegraph(ctx context.Context, message string) (string, error) {
	telegraph.SetSocksDialer(sockAddr)

	if account.AccessToken == "" {
		ac, err := telegraph.CreateAccount(account)
		if err != nil {
			return "", err
		}
		account.AccessToken = ac.AccessToken
	}

	if message == "" {
		message = template
	}

	content, err := telegraph.ContentFormat(message)
	if err != nil {
		return "", err
	}

	// Create new Telegraph page
	pageData := telegraph.Page{
		Title:   "Techcats List", // required
		Content: content,         // required

		// Not necessarily, but, hey, it's just an example.
		AuthorName: account.AuthorName, // optional
		AuthorURL:  account.AuthorURL,  // optional
	}

	page, err := account.CreatePage(pageData, false)
	log.Println(page.URL)
	// todo: send to telegram

	tg.Send(tgbotapi.NewMessage(tg.chatID, page.URL))

	return page.URL, err
}

var template = `
<figure>
	<img src="/file/6a5b15e7eb4d7329ca7af.jpg"/>
</figure>
<p><i>Hello</i>, my name is <b>Tech Cats</b>, <u>look at me</u>!</p>
<figure>
	<iframe src="https://youtu.be/fzQ6gRAEoy0"></iframe>
	<figcaption>
		Yes, you can embed youtube, vimeo and twitter widgets too!
	</figcaption>
</figure>
`
