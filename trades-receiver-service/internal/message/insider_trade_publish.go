package message

import (
	"context"
	"fmt"
	"time"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/rabbitmq"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/config"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/gofrs/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InsiderTradePublisher struct {
	rmq    *rabbitmq.RabbitMQ
	rmqCfg config.RabbitMQ
}

func New(rabbitMQ *rabbitmq.RabbitMQ, rmqCfg config.RabbitMQ) (*InsiderTradePublisher, error) {

	err := rabbitMQ.Channel.ExchangeDeclare(
		rmqCfg.Exchange,
		amqp.ExchangeFanout,
		rmqCfg.Durable,
		false,
		false,
		false,
		nil)

	if err != nil {
		return nil, fmt.Errorf("message New: %w", err)
	}

	q, err := rabbitMQ.Channel.QueueDeclare(
		rmqCfg.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("message New: %w", err)
	}

	err = rabbitMQ.Channel.QueueBind(
		q.Name,
		rmqCfg.RoutingKey,
		rmqCfg.Exchange,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("message New: %w", err)
	}

	return &InsiderTradePublisher{rmq: rabbitMQ, rmqCfg: rmqCfg}, nil
}

func (p *InsiderTradePublisher) PublishTrade(ctx context.Context, trade *entity.Trade) error {

	insiderProto := api.Insider{
		ID:   trade.Ins.ID,
		Cik:  trade.Ins.Cik,
		Name: trade.Ins.Name,
	}

	companyProto := api.Company{
		ID:     trade.Cmp.ID,
		Cik:    trade.Cmp.Cik,
		Name:   trade.Cmp.Name,
		Ticker: trade.Cmp.Ticker,
	}

	secFilingProto := api.SecFiling{
		ID:              trade.SecF.ID,
		FilingType:      trade.SecF.FilingType,
		URL:             trade.SecF.URL,
		InsiderID:       trade.SecF.InsiderID,
		OfficerPosition: trade.SecF.OfficerPosition,
		ReportedOn:      trade.SecF.ReportedOn,
	}

	var securityTransactionHoldingsProto []*api.SecurityTransactionHoldings

	for _, sth := range trade.Sth {
		pricePerSecurity, _ := sth.PricePerSecurity.Float64()

		sthProto := &api.SecurityTransactionHoldings{
			ID:                                sth.ID,
			TransactionID:                     sth.TransactionID,
			SecFilingsID:                      sth.SecFilingsID,
			QuantityOwnedFollowingTransaction: sth.QuantityOwnedFollowingTransaction,
			SecurityTitle:                     sth.SecurityTitle,
			SecurityType:                      sth.SecurityType,
			Quantity:                          sth.Quantity,
			PricePerSecurity:                  pricePerSecurity,
			TransactionDate:                   sth.TransactionDate,
			TransactionCode:                   sth.TransactionCode,
		}

		securityTransactionHoldingsProto = append(securityTransactionHoldingsProto, sthProto)
	}

	transactionProto := api.Transaction{
		ID:                  trade.Trs.ID,
		SecFilingsID:        trade.Trs.SecFilingsID,
		TransactionTypeName: trade.Trs.TransactionTypeName,
		AveragePrice:        trade.Trs.AveragePrice,
		TotalShares:         trade.Trs.TotalShares,
		TotalValue:          trade.Trs.TotalValue,
		CreatedAt:           timestamppb.New(trade.Trs.CreatedAt),
	}

	tradeProto := api.Trade{
		Ins:  &insiderProto,
		Cmp:  &companyProto,
		SecF: &secFilingProto,
		Trs:  &transactionProto,
		Sth:  securityTransactionHoldingsProto,
	}

	encodedTrade, err := proto.Marshal(&tradeProto)
	if err != nil {
		return fmt.Errorf("PublishTrade: failed to marshal trade into proto %w", err)
	}

	if err := p.publish(encodedTrade); err != nil {
		return fmt.Errorf("PublishTrade: %w", err)
	}

	return nil
}

func (p *InsiderTradePublisher) publish(body []byte) error {
	msgID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("publish: failed to generate a new uuid: %w", err)
	}

	msg := amqp.Publishing{
		MessageId:       msgID.String(),
		Timestamp:       time.Now(),
		DeliveryMode:    amqp.Persistent,
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            body,
	}

	if err := p.rmq.Channel.Publish(
		p.rmqCfg.Exchange,
		p.rmqCfg.RoutingKey,
		false,
		false,
		msg,
	); err != nil {
		return fmt.Errorf("publish: failed to publish: %w", err)
	}

	return nil
}
