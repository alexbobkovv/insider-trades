package service

import (
	"fmt"

	"github.com/alexbobkovv/insider-trades/api"
)

type notificationService struct {
	api TelegramAPI
}

func New(tgAPI TelegramAPI) *notificationService {
	return &notificationService{api: tgAPI}
}

func (n *notificationService) ProcessTrade(trade *api.Trade) error {
	const minTotalValue = 10000
	if trade.Trs.TotalValue < minTotalValue {
		return nil
	}

	if err := n.api.SendTrade(trade); err != nil {
		return fmt.Errorf("service: ProcessTrade: %v", err)
	}

	return nil
}
