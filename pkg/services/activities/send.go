package activities

import (
	"context"

	"go.uber.org/zap"
)

func SendMessage(ctx context.Context, message interface{}) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("message", message)
	return nil
}
