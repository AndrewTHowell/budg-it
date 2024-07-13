package dbconvert

import (
	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
)

func ToPayees(dbPayees ...*db.Payee) []*budgit.Payee {
	payees := make([]*budgit.Payee, 0, len(dbPayees))
	for _, dbPayee := range dbPayees {
		payees = append(payees, toPayee(dbPayee))
	}
	return payees
}

func toPayee(payee *db.Payee) *budgit.Payee {
	return &budgit.Payee{
		ID:   payee.ID.String,
		Name: payee.Name.String,
	}
}

func FromPayees(payees ...*budgit.Payee) []*db.Payee {
	dbPayees := make([]*db.Payee, 0, len(payees))
	for _, payee := range payees {
		dbPayees = append(dbPayees, fromPayee(payee))
	}
	return dbPayees
}

func fromPayee(payee *budgit.Payee) *db.Payee {
	return &db.Payee{
		ID:   toText(payee.ID),
		Name: toText(payee.Name),
	}
}
