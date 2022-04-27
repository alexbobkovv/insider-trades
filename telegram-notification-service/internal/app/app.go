package app

import (
	"log"

	"github.com/alexbobkovv/insider-trades/pkg/logger"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"
)

func Run(cfg *config.Config) {
	l, err := logger.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.Filepath)
	if err != nil {
		log.Println("app: failed to initialize logger")
	}
	l.Info("app: logger initialized")

}
