package repository

import (
	"context"
	"fmt"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/postgresql"
	"github.com/jackc/pgx/v4"
)

type InsiderTradeRepo struct {
	*postgresql.Postgres
	l *logger.Logger
}

func New(db *postgresql.Postgres, logger *logger.Logger) *InsiderTradeRepo {
	return &InsiderTradeRepo{db, logger}
}

func (r *InsiderTradeRepo) StoreTrade(ctx context.Context, trade *entity.Trade) error {

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			r.l.Error("transaction failed: ", err)
		}
	}()

	const insiderInsertQuery = `
							WITH insiderID AS (
								INSERT INTO insider (cik, name)
									VALUES ($1, $2)
									ON CONFLICT (cik) DO NOTHING
									RETURNING id
							)
							SELECT *
							FROM insiderID
							UNION
							SELECT id
							FROM insider
							WHERE cik=$1 
								AND name=$2`

	err = tx.QueryRow(ctx, insiderInsertQuery, trade.Ins.Cik, trade.Ins.Name).Scan(&trade.Ins.ID)

	if err != nil {
		return err
	}

	const companyInsertQuery = `
							WITH companyID AS (
								INSERT INTO company (cik, name, ticker)
									VALUES ($1, $2, $3)
									ON CONFLICT (cik) DO NOTHING
									RETURNING id
							)
							SELECT *
							FROM companyID
							UNION
							SELECT id
							FROM company
							WHERE cik = $1
							  AND name = $2
							  AND ticker = $3`

	err = tx.QueryRow(ctx, companyInsertQuery, trade.Cmp.Cik,
		trade.Cmp.Name, trade.Cmp.Ticker).Scan(&trade.Cmp.ID)

	if err != nil {
		return err
	}

	trade.SecF.InsiderID = trade.Ins.ID
	trade.SecF.CompanyID = trade.Cmp.ID

	const secFilingInsertQuery = `
							INSERT INTO sec_filings (filing_type, url, insider_id, company_id, officer_position, reported_on)
							VALUES ($1, $2, $3, $4, $5, $6)
							RETURNING id`

	err = tx.QueryRow(ctx, secFilingInsertQuery, trade.SecF.FilingType,
		trade.SecF.URL, trade.SecF.InsiderID, trade.SecF.CompanyID, trade.SecF.OfficerPosition,
		trade.SecF.ReportedOn).Scan(&trade.SecF.ID)

	if err != nil {
		return err
	}

	const transactionInsertQuery = `
							INSERT INTO transaction 
							(sec_filings_id, transaction_type_name, average_price, total_shares, total_value)
							VALUES ($1, $2, $3, $4, $5)
							RETURNING id`

	err = tx.QueryRow(ctx, transactionInsertQuery, trade.SecF.ID,
		trade.Trs.TransactionTypeName, trade.Trs.AveragePrice, trade.Trs.TotalShares,
		trade.Trs.TotalValue).Scan(&trade.Trs.ID)

	if err != nil {
		return err
	}

	const securityTransactionHoldingsInsertQuery = `
							INSERT INTO security_transaction_holdings 
							(transaction_id, sec_filings_id, quantity_owned_following_transaction,
							security_title, security_type, quantity, price_per_security,
							transaction_date, transaction_code)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
							RETURNING id`

	cthBatch := &pgx.Batch{}

	for _, sth := range trade.Sth {
		cthBatch.Queue(securityTransactionHoldingsInsertQuery, trade.Trs.ID, trade.SecF.ID,
			sth.QuantityOwnedFollowingTransaction, sth.SecurityTitle, sth.SecurityType,
			sth.Quantity, sth.PricePerSecurity, sth.TransactionDate, sth.TransactionCode)
	}

	br := tx.SendBatch(ctx, cthBatch)

	var sthIDs []string

	for i := 0; i < len(trade.Sth); i++ {
		var sthId string

		err = br.QueryRow().Scan(&sthId)
		if err != nil {
			return err
		}

		sthIDs = append(sthIDs, sthId)
	}

	if err = br.Close(); err != nil {
		return err
	}

	if len(sthIDs) != len(trade.Sth) {
		return fmt.Errorf("the length of the ids array from db doesn't match trade.Sth")
	}

	for i := 0; i < len(trade.Sth); i++ {
		trade.Sth[i].ID = sthIDs[i]
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *InsiderTradeRepo) GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	return []*entity.Transaction{{}}, nil
}
