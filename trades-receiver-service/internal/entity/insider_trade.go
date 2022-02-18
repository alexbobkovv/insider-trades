package entity

type Insider struct {
	ID   string
	Cik  int
	Name string
}

type Company struct {
	ID     string
	Cik    int
	Name   string
	Ticker string
}

type SecFilings struct {
	ID              string
	FilingType      int
	URL             string
	InsiderID       string
	CompanyID       string
	OfficerPosition string
	ReportedOn      string
}

type Transaction struct {
	ID                  string
	SecFilingsID        string
	TransactionTypeName string
	AveragePrice        float64
	TotalShares         int
	TotalValue          float64
	CreatedAt           string
}

type SecurityTransactionHoldings struct {
	ID                                string
	TransactionID                     string
	SecFilingsID                      string
	QuantityOwnedFollowingTransaction float64
	SecurityTitle                     string
	SecurityType                      int
	Quantity                          int
	PricePerSecurity                  float64
	TransactionDate                   string
	TransactionCode                   int
	CreatedAt                         string
}
