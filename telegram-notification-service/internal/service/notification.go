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

// ProcessTrade checks constraints and publishes trade to telegram channels
func (n *notificationService) ProcessTrade(trade *api.Trade) error {
	// Ignore trades with total value less than 10 000$ for telegram channel just to keep feed cleaner
	const minTotalValue = 10000
	if trade.Trs.TotalValue < minTotalValue {
		return nil
	}

	// Send(publish) trade to telegram channel
	if err := n.api.SendTrade(trade); err != nil {
		return fmt.Errorf("service: ProcessTrade: %v", err)
	}

	return nil
}
