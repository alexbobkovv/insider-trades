package service

import (
	"github.com/alexbobkovv/insider-trades/api"
)

type (
	Service interface {
		ProcessTrade(trade *api.Trade) error
	}

	TelegramAPI interface {
		SendTrade(trade *api.Trade) error
	}
)
