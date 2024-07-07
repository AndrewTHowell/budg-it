package budgit

import (
	date "cloud.google.com/go/civil"
	"github.com/google/uuid"
)

// Transaction is an Transaction, unique only within a given Budget.
type Transaction struct {
	ID              string
	Date            date.Date
	AccountID       string
	PayeeID         string
	IsPayeeInternal bool
	CategoryID      string
	Amount          BalanceAmount
	Cleared         bool
}

// NewTransaction returns an Transaction in the given Budget.
func NewTransaction(accountID string, date date.Date, payeeID string, isPayeeInternal bool, categoryID string, amount BalanceAmount, cleared bool) *Transaction {
	return &Transaction{
		ID:              uuid.New().String(),
		AccountID:       accountID,
		Date:            date,
		PayeeID:         payeeID,
		IsPayeeInternal: isPayeeInternal,
		CategoryID:      categoryID,
		Amount:          amount,
		Cleared:         cleared,
	}
}

func (t Transaction) GetID() string {
	return t.ID
}

func (t Transaction) Mirror() *Transaction {
	return &Transaction{
		ID:              uuid.New().String(),
		AccountID:       t.PayeeID,
		Date:            t.Date,
		PayeeID:         t.AccountID,
		IsPayeeInternal: t.IsPayeeInternal,
		CategoryID:      t.CategoryID,
		Amount:          -t.Amount,
		Cleared:         t.Cleared,
	}
}
