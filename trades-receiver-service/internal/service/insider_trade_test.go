package service

import (
	"context"
	"testing"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type test struct {
	name string
	mock func()
	// res  interface{}
	err error
}

func TestReceive(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockInsiderTradeRepo(ctrl)
	publisher := NewMockInsiderTradePublisher(ctrl)

	tradeService := New(repo, publisher)

	ins := &entity.Insider{
		Cik:  1878229,
		Name: "MONROE WILLIAM",
	}
	cmp := &entity.Company{
		Cik:    1537028,
		Name:   "Independence Contract Drilling, Inc.",
		Ticker: "ICD",
	}
	fType := 0
	secF := &entity.SecFiling{
		FilingType: &fType,
		URL:        "https://www.sec.gov/Archives/edgar/data/1537028/000187822922000005/0001878229-22-000005-index.htm",
		ReportedOn: "2022-02-04T00:00:00",
	}

	quantity := 2048000.0
	secType := 0
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

	trade := &entity.Trade{
		Ins: ins, Cmp: cmp, SecF: secF, Sth: sth,
	}

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().StoreTrade(context.Background(), trade).Return(nil)
			},
			// res: entity.Trade{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()
			err := tradeService.Receive(context.Background(), &entity.Trade{})

			require.ErrorIs(t, err, tc.err)
		})
	}

}
