package budgit

import (
	"time"
)

// Transaction is an Transaction, unique only within a given Budget.
type Transaction struct {
	ID              string
	EffectiveDate   time.Time
	AccountID       string
	PayeeID         string
	IsPayeeInternal bool
	CategoryID      string
	Amount          BalanceAmount
	Cleared         bool
}

// Mirror mirrors the transaction, by returning another with the same fields but:
//   - ID is replaced with given ID
//   - Account and Payee IDs are swapped
//   - Amount is negated
func (t Transaction) Mirror(id string) *Transaction {
	return &Transaction{
		ID:              id,
		AccountID:       t.PayeeID,
		EffectiveDate:   t.EffectiveDate,
		PayeeID:         t.AccountID,
		IsPayeeInternal: t.IsPayeeInternal,
		CategoryID:      t.CategoryID,
		Amount:          -t.Amount,
		Cleared:         t.Cleared,
	}
}
