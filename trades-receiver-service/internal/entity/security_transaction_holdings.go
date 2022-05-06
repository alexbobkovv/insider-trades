package entity

import "math/big"

type SecurityTransactionHoldings struct {
	ID                                string
	TransactionID                     *string
	SecFilingsID                      string
	QuantityOwnedFollowingTransaction *float64
	SecurityTitle                     string
	SecurityType                      *int32
	Quantity                          int64
	PricePerSecurity                  big.Rat
	TransactionDate                   string
	TransactionCode                   int32
}
