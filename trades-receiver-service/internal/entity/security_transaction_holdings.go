package entity

type SecurityTransactionHoldings struct {
	ID                                string
	TransactionID                     *string
	SecFilingsID                      string
	QuantityOwnedFollowingTransaction *float64
	SecurityTitle                     string
	SecurityType                      *int
	Quantity                          int
	PricePerSecurity                  float64
	TransactionDate                   string
	TransactionCode                   int
}
