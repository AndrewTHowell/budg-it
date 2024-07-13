package budgit

import (
	"time"

	"github.com/google/uuid"
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

// NewTransaction returns an Transaction in the given Budget.
func NewTransaction(accountID string, date time.Time, payeeID string, isPayeeInternal bool, categoryID string, amount BalanceAmount, cleared bool) *Transaction {
	return &Transaction{
		ID:              uuid.New().String(),
		AccountID:       accountID,
		EffectiveDate:   date,
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
		EffectiveDate:   t.EffectiveDate,
		PayeeID:         t.AccountID,
		IsPayeeInternal: t.IsPayeeInternal,
		CategoryID:      t.CategoryID,
		Amount:          -t.Amount,
		Cleared:         t.Cleared,
	}
}
