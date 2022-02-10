package main

import (
	"insidertradesreceiver/config"
	"insidertradesreceiver/internal/app"
	"insidertradesreceiver/pkg/logger"
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