package main

import (
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/config"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/app"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"
)

func main() {
	l := logger.New()
	l.Info("Logger initialized")

	cfg, err := config.New()
	if err != nil {
		l.Fatal("Failed to load server config: %v", err)
	}

	app.Run(cfg, l)
}
