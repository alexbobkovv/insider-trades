package telegram

import (
	"fmt"
	"time"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	printer := message.NewPrinter(language.English)
	reportedOn, err := time.Parse(time.RFC3339, trade.SecF.ReportedOn)
	if err != nil {
		return fmt.Errorf("tgbot_api: SendTrade: failed to parse SecF reportedOn: %w", err)
	}

	msgText := fmt.Sprintf(
		"Ticker: #%s\n"+
			"Company: <b>%s</b>\n"+
			"Insider: <b>%s</b>\n"+
			"Type: <b>%s</b>\n"+
			"Total shares: <b>%v</b>\n"+
			"Average price: <b>%s$</b>\n"+
			"Total value: <b>%s$</b>\n"+
			"<a href=\"%s\">SEC form</a>\n"+
			"Reported on: <b>%v</b>",
		trade.Cmp.Ticker,
		trade.Cmp.Name,
		trade.Ins.Name,
		trade.Trs.TransactionTypeName,
		trade.Trs.TotalShares,
		printer.Sprintf("%.3f", trade.Trs.AveragePrice),
		printer.Sprintf("%.3f", trade.Trs.TotalValue),
		trade.SecF.URL,
		reportedOn.Format(time.RFC1123),
	)

	msg := tgbotapi.NewMessage(t.cfg.ChannelID, msgText)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = t.api.Send(msg)
	if err != nil {
		return fmt.Errorf("telegram: SendTrade: failed to send message to telegram channel: %w", err)
	}

	return nil
}
