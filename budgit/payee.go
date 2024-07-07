package budgit

import (
	"github.com/google/uuid"
)

// Payee is an external party involved in a payee.
type Payee struct {
	ID   string
	Name string
}

// NewPayee returns an Payee.
func NewPayee(name string) *Payee {
	return &Payee{
		ID:   uuid.New().String(),
		Name: name,
	}
}

func (p Payee) GetID() string {
	return p.ID
}
