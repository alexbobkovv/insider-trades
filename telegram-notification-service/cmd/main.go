package main

import (
	"log"

	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("main: failed to load config: %v", err)
	}

	app.Run(cfg)
}
