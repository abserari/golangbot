// Binary gotdecho provides example of Telegram echo bot.
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/proxy"
	"golang.org/x/xerrors"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/transport"
	"github.com/spf13/viper"
)

func run(ctx context.Context) error {
	logger, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.InfoLevel))
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
	// from https://my.telegram.org/apps
	appID := viper.GetInt("TgAppID") // integer value from "App api_id" field
	appHash := viper.GetString("TgAppHash")

	// Setting up session storage.
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	sessionDir := filepath.Join(home, ".td")
	if err := os.MkdirAll(sessionDir, 0600); err != nil {
		return err
	}

	dispatcher := tg.NewUpdateDispatcher()
	// Creating connection.
	dialCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	sock5, err := proxy.SOCKS5("tcp", "106.53.131.105:65431", &proxy.Auth{
		User:     "Kassulke8264",
		Password: "wFKo1z8xOr",
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

				return client.SendMessage(ctx, &tg.MessagesSendMessageRequest{
					Message: m.Message,
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
		return xerrors.Errorf("failed to connect: %w", err)
	}
	logger.Info("Connected")

	self, err := client.Self(ctx)
	if err != nil || !self.Bot {
		if err := client.AuthBot(dialCtx, os.Getenv("TELEGRAM_TECHCATS_BOT_TOKEN")); err != nil {
			return xerrors.Errorf("failed to perform bot login: %w", err)
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
		return xerrors.Errorf("failed to get state: %w", err)
	}
	logger.Sugar().Infof("Got state: %+v", state)

	// Reading updates until SIGTERM.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logger.Info("Shutting down")
	if err := client.Close(); err != nil {
		return err
	}
	logger.Info("Graceful shutdown completed")
	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(2)
	}
}
