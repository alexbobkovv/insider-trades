package main

import (
	"log"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/config"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("main: failed to load server config: %v", err)
	}

	app.Run(cfg)
}
