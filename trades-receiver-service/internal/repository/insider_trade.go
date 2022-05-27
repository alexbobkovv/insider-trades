package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/postgresql"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		pricePerSecurity, _ := sth.PricePerSecurity.Float64()
		cthBatch.Queue(securityTransactionHoldingsInsertQuery, trade.Trs.ID, trade.SecF.ID,
			sth.QuantityOwnedFollowingTransaction, sth.SecurityTitle, sth.SecurityType,
			sth.Quantity, pricePerSecurity, sth.TransactionDate, sth.TransactionCode)
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

func (r *InsiderTradeRepo) ListTransactions(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*entity.Transaction, *cursor.Cursor, error) {
	const methodName = "(r *InsiderTradeRepo) ListTrades"

	var rows pgx.Rows
	if reqCursor.IsEmpty() {
		const transactionSelectQuery = `
			SELECT id, sec_filings_id, transaction_type_name, average_price, total_shares, total_value, created_at
			FROM transaction
			ORDER BY created_at DESC
			LIMIT $1`

		var err error
		rows, err = r.Pool.Query(ctx, transactionSelectQuery, limit)

		if err != nil {
			return nil, cursor.NewEmpty(), fmt.Errorf("%v: %w", methodName, err)
		}

	} else {
		decodedCursor := reqCursor.GetDecoded()

		const transactionSelectCursorQuery = `
			SELECT id, sec_filings_id, transaction_type_name, average_price, total_shares, total_value, created_at
			FROM transaction
			WHERE created_at < $1 :: timestamptz
			ORDER BY created_at DESC
			LIMIT $2`

		var err error
		rows, err = r.Pool.Query(ctx, transactionSelectCursorQuery, *decodedCursor, limit)

		if err != nil {
			return nil, cursor.NewEmpty(), fmt.Errorf("%v: %w", methodName, err)
		}
	}

	var transactions []*entity.Transaction
	var averagePrice, totalValue float64

	for rows.Next() {
		var transaction entity.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.SecFilingsID,
			&transaction.TransactionTypeName,
			&averagePrice,
			&transaction.TotalShares,
			&totalValue,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, cursor.NewEmpty(), fmt.Errorf("%v: %w", methodName, err)
		}
		transaction.AveragePrice = averagePrice
		transaction.TotalValue = totalValue

		transactions = append(transactions, &transaction)
	}

	if len(transactions) > 0 {
		cursorTimestamp := transactions[len(transactions)-1].CreatedAt
		nextCursor := cursor.NewFromTime(&cursorTimestamp)
		return transactions, nextCursor, nil
	}

	return transactions, cursor.NewEmpty(), nil
}

func (r *InsiderTradeRepo) ListViews(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*api.TradeViewResponse, *cursor.Cursor, error) {
	const methodName = "(r *InsiderTradeRepo) ListViews"

	var rows pgx.Rows
	if reqCursor.IsEmpty() {
		const tradeViewsSelectQuery = `
			SELECT id,
				   sec_filings_id,
				   transaction_type_name,
				   average_price,
				   total_shares,
				   total_value,
				   created_at,
				   url,
				   insider_id,
				   company_id,
				   officer_position,
				   reported_on,
				   insider_cik,
				   insider_name,
				   company_cik,
				   company_name,
				   ticker
			FROM trades_matview
			LIMIT $1
			`

		var err error
		rows, err = r.Pool.Query(ctx, tradeViewsSelectQuery, limit)

		if err != nil {
			return nil, nil, fmt.Errorf("%v: %w", methodName, err)
		}

	} else {

		const tradeViewsSelectCursorQuery = `
			SELECT id,
				sec_filings_id,
				transaction_type_name,
				average_price,
				total_shares,
				total_value,
				created_at,
				url,
				insider_id,
				company_id,
				officer_position,
				reported_on,
				insider_cik,
				insider_name,
				company_cik,
				company_name,
				ticker
			FROM trades_matview
			WHERE created_at < $1 :: timestamptz
			LIMIT $2
			`

		var err error
		rows, err = r.Pool.Query(ctx, tradeViewsSelectCursorQuery, *reqCursor.GetDecoded(), limit)

		if err != nil {
			return nil, nil, fmt.Errorf("%v: %w", methodName, err)
		}
	}

	var tradeViews []*api.TradeViewResponse

	for rows.Next() {
		var tradeView api.TradeViewResponse

		var reportedOn time.Time
		var createdAt time.Time
		err := rows.Scan(
			&tradeView.ID,
			&tradeView.SecFilingsID,
			&tradeView.TransactionTypeName,
			&tradeView.AveragePrice,
			&tradeView.TotalShares,
			&tradeView.TotalValue,
			&createdAt,
			&tradeView.URL,
			&tradeView.InsiderID,
			&tradeView.CompanyID,
			&tradeView.OfficerPosition,
			&reportedOn,
			&tradeView.InsiderCik,
			&tradeView.InsiderName,
			&tradeView.CompanyCik,
			&tradeView.CompanyName,
			&tradeView.CompanyTicker,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: failed to scan tradeView: %w", methodName, err)
		}

		const dateLayout = "01-02-2006"
		tradeView.ReportedOn = reportedOn.Format(dateLayout)
		tradeView.CreatedAt = timestamppb.New(createdAt)

		tradeViews = append(tradeViews, &tradeView)
	}

	if len(tradeViews) > 0 {
		cursorTimestamp := tradeViews[len(tradeViews)-1].CreatedAt.AsTime()
		nextCursor := cursor.NewFromTime(&cursorTimestamp)

		return tradeViews, nextCursor, nil
	}

	return tradeViews, cursor.NewEmpty(), nil

}

func (r *InsiderTradeRepo) RefreshTradeMatView(ctx context.Context) error {
	const methodName = "(r *InsiderTradeRepo) RefreshTradeMatView"

	const refreshTradeMatViewQuery = `
		REFRESH MATERIALIZED VIEW trades_matview
        `

	_, err := r.Pool.Exec(ctx, refreshTradeMatViewQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", methodName, err)
	}

	return nil
}
