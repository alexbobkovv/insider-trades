package service_test

import (
	"context"
	"testing"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestReceive(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockInsiderTradeRepo(ctrl)
	publisher := NewMockInsiderTradePublisher(ctrl)

	tradeService := service.New(repo, publisher)

	ins := &entity.Insider{
		Cik:  1878229,
		Name: "MONROE WILLIAM",
	}
	cmp := &entity.Company{
		Cik:    1537028,
		Name:   "Independence Contract Drilling, Inc.",
		Ticker: "ICD",
	}
	var fType int64
	secF := &entity.SecFiling{
		FilingType: &fType,
		URL:        "https://www.sec.gov/Archives/edgar/data/1537028/000187822922000005/0001878229-22-000005-index.htm",
		ReportedOn: "2022-02-04T00:00:00",
	}

	quantity := 2048000.0
	var secType int32
	sth := []*entity.SecurityTransactionHoldings{
		{
			QuantityOwnedFollowingTransaction: &quantity,
			SecurityTitle:                     "Common Stock",
			SecurityType:                      &secType,
			Quantity:                          18000.0,
			PricePerSecurity:                  3.17,
			TransactionDate:                   "2022-02-04T00:00:00",
			TransactionCode:                   0,
		},
	}

	validTrade := &entity.Trade{
		Ins: ins, Cmp: cmp, SecF: secF, Sth: sth,
	}

	tests := []struct {
		name string
		args *entity.Trade
		mock func()
		// res  interface{}
		// err error
		wantErr bool
	}{
		{
			name:    "empty trade",
			args:    &entity.Trade{},
			mock:    nil,
			wantErr: true,
		},
		{
			name: "valid trade",
			args: validTrade,
			mock: func() {
				repo.EXPECT().StoreTrade(context.Background(), validTrade).Return(nil)
				publisher.EXPECT().PublishTrade(context.Background(), validTrade).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.mock != nil {
				tc.mock()
			}

			err := tradeService.Receive(context.Background(), tc.args)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.ErrorIs(t, err, nil)
			}
		})
	}

}
