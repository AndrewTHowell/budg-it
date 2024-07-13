package convert

import (
	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
)

func ToPayee(payee *db.Payee) *budgit.Payee {
	return &budgit.Payee{
		ID:   payee.ID.String,
		Name: payee.Name.String,
	}
}

func FromPayee(payee *budgit.Payee) *db.Payee {
	return &db.Payee{
		ID:   toText(payee.ID),
		Name: toText(payee.Name),
	}
}
