package app

import (
	"log"

	"github.com/alexbobkovv/insider-trades/pkg/rabbitmq"
	"github.com/alexbobkovv/insider-trades/pkg/zap"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run(cfg *config.Config) {
	l, err := zap.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.Filepath)
	if err != nil {
		log.Println("app: failed to initialize zap")
	}
	l.Info("app: zap initialized")

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		l.Fatalf("app: tgbotapi.NewBotAPI: %v", err)
	}
	l.Infof("Authorized on account %s", bot.Self.UserName)

	amqpClient, err := rabbitmq.NewClient(cfg.AmqpURL)
	if err != nil {
		l.Fatalf("app: failed to connect to RabbitMQ")
	}

	msgs, err := amqpClient.Consume(
		"telegram_channel_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		l.Fatalf("failed to register a consumer")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			l.Info("%s", d.Body)
		}
	}()

	l.Info("waiting for messages..")
	<-forever
	// tradesChan := tgbotapi.Chat{
	// 	ID:   cfg.ChannelID,
	// 	Type: "channel",
	// }
	//
	// msg := tgbotapi.NewMessage(tradesChan.ID, ":))")
	// bot.Send(msg)
}
