package db

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
)

func (db *DB) InsertPayees(ctx context.Context, payees ...*budgit.Payee) error {
	if err := insert(&db.payees, payees...); err != nil {
		return fmt.Errorf("inserting %d payees: %w", len(payees), err)
	}
	return nil
}

func (db *DB) SelectPayeesByID(ctx context.Context, payeeIDs ...string) (map[string]*budgit.Payee, error) {
	return selectByIDs(db.payees, payeeIDs), nil
}

func (db *DB) SelectPayeesByName(ctx context.Context, payeeNames ...string) (map[string]*budgit.Payee, error) {
	targetName := boolMap(payeeNames)
	payees := make(map[string]*budgit.Payee, len(payeeNames))
	for _, payee := range db.payees {
		if targetName[payee.Name] {
			payees[payee.Name] = payee
		}
	}
	return payees, nil
}
