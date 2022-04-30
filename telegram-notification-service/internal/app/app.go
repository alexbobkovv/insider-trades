package app

import (
	"log"

	"github.com/alexbobkovv/insider-trades/pkg/rabbitmq"
	"github.com/alexbobkovv/insider-trades/pkg/zap"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"
)

func Run(cfg *config.Config) {
	l, err := zap.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.Filepath)
	if err != nil {
		log.Println("app: failed to initialize zap")
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

	msgs, err := rmq.Channel.Consume(
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
			l.Infof("%s", d.Body)
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
