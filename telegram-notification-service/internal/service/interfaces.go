package service

import "github.com/alexbobkovv/insider-trades/telegram-notification-service/internal/entity"

type TelegramAPI interface {
	SendTrade(trade *entity.Trade) error
}
