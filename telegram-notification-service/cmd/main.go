package main

import (
	"log"
	"os"

	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/internal/app"
)

func main() {

	cfg, err := config.New(os.Getenv("CONFIG_PATH"), os.Getenv("CONFIG_NAME"))
	if err != nil {
		log.Fatalf("main: failed to load config: %v", err)
	}

	app.Run(cfg)
}
