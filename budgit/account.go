package budgit

import (
	"github.com/google/uuid"
)

// Account is an Account, unique only within a given Budget.
type Account struct {
	ID                string
	Name              string
	ExternalAccountID string
	Balance           Balance
}

// NewAccount returns an Account in the given Budget.
func NewAccount(name, externalAccountID string, balance Balance) *Account {
	return &Account{
		ID:                uuid.New().String(),
		Name:              name,
		ExternalAccountID: externalAccountID,
		Balance:           balance,
	}
}

func (a Account) GetID() string {
	return a.ID
}

// ExternalAccount is an Account representing some real, external Account.
type ExternalAccount struct {
	ID                 string
	ExternalProviderID string
	ExternalID         string
	Name               string
	Balance            Balance
}

// NewExternalAccount returns an ExternalAccount.
func NewExternalAccount(externalProviderID, externalID, name string, balance Balance) *ExternalAccount {
	return &ExternalAccount{
		ID:                 uuid.New().String(),
		ExternalProviderID: externalProviderID,
		ExternalID:         externalID,
		Name:               name,
		Balance:            balance,
	}
}

func (a ExternalAccount) GetID() string {
	return a.ID
}
