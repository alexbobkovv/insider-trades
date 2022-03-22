package repository

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/postgresql"
	"github.com/jackc/pgx/v4"
)

type InsiderTradeRepo struct {
	*postgresql.Postgres
}

func New(db *postgresql.Postgres) *InsiderTradeRepo {
	return &InsiderTradeRepo{db}
}

func (r *InsiderTradeRepo) StoreTrade(ctx context.Context, trade *entity.Trade) (err error) {
	const methodName = "(r *InsiderTradeRepo) StoreTrade"

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%v: failed to begin transaction: %w", methodName, err)
	}
	defer func() {
		if tempErr := tx.Rollback(ctx); tempErr != nil && tempErr != pgx.ErrTxClosed {
			err = tempErr
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
		return fmt.Errorf("%v: failed to query insert insider: %w", methodName, err)
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
		return fmt.Errorf("%v: failed to query insert company: %w", methodName, err)
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
		return fmt.Errorf("%v: failed to query insert sec_filings: %w", methodName, err)
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
		return fmt.Errorf("%v: failed to query insert transaction: %w", methodName, err)
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
			return fmt.Errorf("%v: failed to query batch id: %w", methodName, err)
		}

		sthIDs = append(sthIDs, sthId)
	}

	if err = br.Close(); err != nil {
		return fmt.Errorf("%v: failed to close batch: %w", methodName, err)
	}

	if len(sthIDs) != len(trade.Sth) {
		return fmt.Errorf("%v: the length of the ids array from db doesn't match trade.Sth", methodName)
	}

	for i := 0; i < len(trade.Sth); i++ {
		trade.Sth[i].ID = sthIDs[i]
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%v: failed to commit transaction: %w", methodName, err)
	}

	return nil
}

func (r *InsiderTradeRepo) decodeTimestampCursor(encodedCursor string) (*time.Time, error) {
	b, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return nil, fmt.Errorf("decodeTimestampCursor: failed to decode cursor: %w", err)
	}

	timestamp, err := time.Parse(time.RFC3339Nano, string(b))
	if err != nil {
		return nil, fmt.Errorf("decodeTimestampCursor: failed parse timestamp: %w", err)
	}

	return &timestamp, nil
}

// TODO encoder
// func (r *InsiderTradeRepo) encodeTimestampCursor(decodedCursor string) {
//
// }

func (r *InsiderTradeRepo) GetAll(ctx context.Context, cursor string, limit int) ([]*entity.Transaction, error) {
	const methodName = "(r *InsiderTradeRepo) GetAll"

	var decodedCursor *time.Time
	if cursor == "" {
		decodedCursor = &time.Time{}
	} else {
		var err error
		decodedCursor, err = r.decodeTimestampCursor(cursor)

		if err != nil {
			return nil, fmt.Errorf("%v: invalid cursor: %w", methodName, err)
		}
	}

	const transactionSelectQuery = `
		SELECT id, sec_filings_id, transaction_type_name, average_price, total_shares, total_value, created_at
		FROM transaction
		WHERE created_at > $1 :: timestamptz
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.Pool.Query(ctx, transactionSelectQuery, *decodedCursor, limit)

	if err != nil {
		return nil, fmt.Errorf("%v: %w", methodName, err)
	}

	var transactions []*entity.Transaction

	for rows.Next() {
		var transaction *entity.Transaction
		err = rows.Scan(
			&transaction.ID,
			&transaction.SecFilingsID,
			&transaction.TransactionTypeName,
			&transaction.AveragePrice,
			&transaction.TotalShares,
			&transaction.TotalValue,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", methodName, err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
