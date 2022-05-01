package service

import "github.com/alexbobkovv/insider-trades/telegram-notification-service/internal/entity"

type (
	Service interface {
		ProcessTrade(trade *entity.Trade) error
	}

	TelegramAPI interface {
		SendTrade(trade *entity.Trade) error
	}
)
