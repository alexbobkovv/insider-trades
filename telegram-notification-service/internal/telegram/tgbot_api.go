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
	msgText := fmt.Sprintf("Ticker: #%s\nCompany: <b>%s</b>\nInsider: <b>%s</b>\n"+
		"Type: <b>%s</b>\nTotal shares: <b>%v</b>\nAverage price: <b>%.2f$</b>\n"+
		"Total value: <b>%.2f$</b>\n"+
		"<a href=\"%s\">SEC form</a>\n"+
		"Reported on: <b>%v</b>",
		trade.Cmp.Ticker,
		trade.Cmp.Name,
		trade.Ins.Name,
		trade.Trs.TransactionTypeName,
		trade.Trs.TotalShares,
		trade.Trs.AveragePrice,
		trade.Trs.TotalValue,
		trade.SecF.URL,
		trade.SecF.ReportedOn,
	)

	msg := tgbotapi.NewMessage(t.cfg.ChannelID, msgText)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := t.api.Send(msg)
	if err != nil {
		return fmt.Errorf("telegram: SendTrade: failed to send message to telegram channel: %w", err)
	}

	return nil
}
