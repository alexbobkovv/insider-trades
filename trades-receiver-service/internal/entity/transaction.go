package entity

import "time"

type Transaction struct {
	ID                  string
	SecFilingsID        string
	TransactionTypeName string
	AveragePrice        float64
	TotalShares         int64
	TotalValue          float64
	CreatedAt           time.Time
}
