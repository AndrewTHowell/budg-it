package budgit

import "github.com/google/uuid"

// Account is an Account, unique only within a given Budget.
type Account struct {
	ID       uuid.UUID
	BudgetID uuid.UUID
	Name     string
}

// NewAccount returns an Account in the given Budget.
func NewAccount(budgetID uuid.UUID, name string) *Account {
	return &Account{
		ID: uuid.New(),
		// TODO: check budget ID existence.
		BudgetID: budgetID,
		Name:     name,
	}
}
