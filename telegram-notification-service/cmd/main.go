package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexbobkovv/insider-trades/pkg/logger"
	"github.com/alexbobkovv/insider-trades/pkg/rabbitmq"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/internal/controller/amqpconsumer"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/internal/service"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/internal/telegram"
)

func main() {
	cfg, err := config.New(os.Getenv("CONFIG_PATH"), os.Getenv("CONFIG_NAME"))
	if err != nil {
		log.Fatalf("main: failed to load config: %v", err)
	}

	l, err := logger.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.Filepath)
	if err != nil {
		log.Printf("app: failed to initialize zap: %v", err)
	}
	l.Info("app: zap initialized")

	rmq, err := rabbitmq.New(cfg.AmqpURL)
	defer func() {
		err := rmq.Channel.Close()
		if err != nil {
			l.Errorf("app: failed to close amqp channel")
		}
	}()
	if err != nil {
		l.Fatalf("app: failed to connect to RabbitMQ: %v", err)
	}

	tgapi, err := telegram.New(&cfg.Telegram)
	if err != nil {
		l.Fatalf("app: failed to initialize telegram api: %v", err)
	}

	notificationService := service.New(tgapi)

	consumer, err := amqpconsumer.New(rmq, &cfg.RabbitMQ, notificationService, l)
	if err != nil {
		l.Fatalf("app: failed to initialize consumer: %v", err)
	}

	// Starting server with graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		l.Info("app: amqp consumer server is running..")
		if err := consumer.Run(); err != nil {
			l.Fatalf("app: faild to start server: %v", err)
		}
	}()

	// Waiting for shutdown signal
	sysSignal := <-signalChan
	l.Infof("app: got signal %v, shutting down..", sysSignal)

	if err := rmq.Channel.Cancel(cfg.ConsumerName, true); err != nil {
		l.Fatalf("app: failed to cancel consumer: %v", err)
	}

	if err := rmq.Connection.Close(); err != nil {
		l.Fatalf("app: failed to close connection to rabbitMQ: %v", err)
	}

	l.Info("app: successful shutdown")
}
