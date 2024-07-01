// Package budgit is the domain package. It contains the domain objects for budgiting.
package budgit

import "github.com/google/uuid"

// Budget is a single Budget used for tracking finances.
type Budget struct {
	ID       uuid.UUID
	Name     string
	Currency Currency
}

// NewBudget returns a Budget.
func NewBudget(name string, currency Currency) *Budget {
	return &Budget{
		ID:       uuid.New(),
		Name:     name,
		Currency: currency,
	}
}

// Currency denotes a currency.
type Currency string

const (
	// GBP is the currency of the United Kingdom.
	GBP = Currency("GBP")
)
