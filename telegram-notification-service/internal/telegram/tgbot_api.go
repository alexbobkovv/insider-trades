package telegram

import (
	"fmt"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBotAPI struct {
	api *tgbotapi.BotAPI
	cfg *config.Telegram
}

func New(cfg *config.Telegram) (*TgBotAPI, error) {
	a, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("telegram New: tgbotapi.NewBotAPI: %w", err)
	}
	return &TgBotAPI{api: a, cfg: cfg}, nil
}

func (t *TgBotAPI) SendTrade(trade *api.Trade) error {

	msg := tgbotapi.NewMessage(t.cfg.ChannelID, trade.String())
	_, err := t.api.Send(msg)
	if err != nil {
		return fmt.Errorf("telegram: SendTrade: failed to send message to telegram channel: %w", err)
	}

	return nil
}
