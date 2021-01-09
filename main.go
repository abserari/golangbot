package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/abserari/golangbot/pkg/pixiv"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/transport"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/proxy"
	"golang.org/x/xerrors"
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
	// from https://my.telegram.org/apps
	appID := viper.GetInt("TgAppID") // integer value from "App api_id" field
	appHash := viper.GetString("TgAppHash")
	botToken := viper.GetString("TgBotToken")

	username := viper.GetString("Username")
	password := viper.GetString("Password")
	cookies := viper.GetString("Cookies")

	socks5Username := viper.GetString("Socks5.Username")
	socks5Pwd := viper.GetString("Socks5.Password")
	socks5Addr := viper.GetString("Socks5.Address")

	log.Println(socks5Addr)

	logger, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.InfoLevel))
	defer func() { _ = logger.Sync() }()

	config := &pixiv.ClientConfig{
		Username: username,
		Password: password,
		Cookies:  cookies,
		Logger:   logger,
	}

	pixivClient := pixiv.NewClient(config)

	ctx := context.Background()
	// Setting up session storage.
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Error(err.Error())
	}
	sessionDir := filepath.Join(home, ".td")
	if err := os.MkdirAll(sessionDir, 0600); err != nil {
		logger.Error(err.Error())
	}

	dispatcher := tg.NewUpdateDispatcher()
	// Creating connection.
	dialCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	sock5, err := proxy.SOCKS5("tcp", socks5Addr, &proxy.Auth{
		User:     socks5Username,
		Password: socks5Pwd,
	}, proxy.Direct)

	dc := sock5.(interface {
		DialContext(ctx context.Context, network, addr string) (net.Conn, error)
	})

	client := telegram.NewClient(appID, appHash, telegram.Options{
		Logger: logger,
		SessionStorage: &telegram.FileSessionStorage{
			Path: filepath.Join(sessionDir, "session.json"),
		},
		Transport: transport.Intermediate(transport.DialFunc(dc.DialContext)),
		// Transport:     trp,
		UpdateHandler: dispatcher.Handle,
	})

	dispatcher.OnNewMessage(func(ctx tg.UpdateContext, u *tg.UpdateNewMessage) error {
		switch m := u.Message.(type) {
		case *tg.Message:
			switch peer := m.PeerID.(type) {
			case *tg.PeerUser:
				user := ctx.Users[peer.UserID]
				logger.With(
					zap.String("text", m.Message),
					zap.Int("user_id", user.ID),
					zap.String("user_first_name", user.FirstName),
					zap.String("username", user.Username),
				).Info("Got message")

				items, err := pixivClient.Ranking(context.Background(), "", "", time.Now().AddDate(0, 0, -1), 0)
				if err != nil {
					logger.Info(err.Error())
					return err
				}

				return client.SendMessage(ctx, &tg.MessagesSendMessageRequest{
					Message: items[0].Image.Regular,
					Peer: &tg.InputPeerUser{
						UserID:     user.ID,
						AccessHash: user.AccessHash,
					},
				})
			}
		}

		return nil
	})

	if err := client.Connect(ctx); err != nil {
		logger.Error(xerrors.Errorf("failed to connect: %w", err).Error())
	}
	logger.Info("Connected")

	self, err := client.Self(ctx)
	if err != nil || !self.Bot {
		if err := client.AuthBot(dialCtx, botToken); err != nil {
			logger.Error(xerrors.Errorf("failed to perform bot login: %w", err).Error())
		}
		logger.Info("Bot login ok")
	}

	// Using tg.Client for directly calling RPC.
	raw := tg.NewClient(client)

	// Getting state is required to process updates in your code.
	// Currently missed updates are not processed, so only new
	// messages will be handled.
	state, err := raw.UpdatesGetState(ctx)
	if err != nil {
		logger.Error(xerrors.Errorf("failed to get state: %w", err).Error())
	}
	logger.Sugar().Infof("Got state: %+v", state)

	// Reading updates until SIGTERM.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logger.Info("Shutting down")
	if err := client.Close(); err != nil {
		logger.Error(err.Error())
	}
	logger.Info("Graceful shutdown completed")
}
