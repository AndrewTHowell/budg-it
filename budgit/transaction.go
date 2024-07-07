package budgit

import (
	date "cloud.google.com/go/civil"
	"github.com/google/uuid"
)

// Transaction is an Transaction, unique only within a given Budget.
type Transaction struct {
	ID        string
	BudgetID  string
	Date      date.Date
	AccountID string
	// TODO: Either external payee or internal Account.
	PayeeID    string
	CategoryID string
	Amount     int
}

// NewTransaction returns an Transaction in the given Budget.
func NewTransaction(budgetID, accountID string, date date.Date, payeeID, categoryID string, amount int) *Transaction {
	return &Transaction{
		ID: uuid.New().String(),
		// TODO: check budget ID existence.
		BudgetID:   budgetID,
		AccountID:  accountID,
		Date:       date,
		PayeeID:    payeeID,
		CategoryID: categoryID,
		Amount:     amount,
	}
}

func (t Transaction) GetID() string {
	return t.ID
}
