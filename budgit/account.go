package budgit

import "github.com/google/uuid"

// Account is an Account, unique only within a given Budget.
type Account struct {
	ID       string
	BudgetID string
	Name     string
}

// NewAccount returns an Account in the given Budget.
func NewAccount(budgetID, name string) *Account {
	return &Account{
		ID: uuid.New().String(),
		// TODO: check budget ID existence.
		BudgetID: budgetID,
		Name:     name,
	}
}

// ExternalAccount is an Account representing some real, external Account.
type ExternalAccount struct {
	ID               string
	ExternalID       string
	ExternalProvider string
}

// NewExternalAccount returns an ExternalAccount.
func NewExternalAccount(externalID, externalProvider string) *ExternalAccount {
	return &ExternalAccount{
		ID:               uuid.New().String(),
		ExternalID:       externalID,
		ExternalProvider: externalProvider,
	}
}
