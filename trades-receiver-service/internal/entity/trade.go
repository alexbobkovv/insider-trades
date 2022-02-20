package entity

type Trade struct {
	Ins  *Insider
	Cmp  *Company
	SecF *SecFiling
	Trs  *Transaction
	Sth  []*SecurityTransactionHoldings
}
