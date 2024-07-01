// Package budgit is the domain package. It contains the domain objects for budgiting.
package budgit

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidCurrency = fmt.Errorf("given currency is not valid")
)

// Budget is a single Budget used for tracking finances.
type Budget struct {
	ID       string
	Name     string
	Currency string
}

// NewBudget returns a Budget.
func NewBudget(name, currency string) (*Budget, error) {
	return &Budget{
		ID:       uuid.New().String(),
		Name:     name,
		Currency: currency,
	}, nil
}
