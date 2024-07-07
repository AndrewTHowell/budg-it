package db

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
)

func (db *DB) InsertTransactions(ctx context.Context, transactions ...*budgit.Transaction) error {
	if err := insert(&db.transactions, transactions...); err != nil {
		return fmt.Errorf("inserting %d transactions: %w", len(transactions), err)
	}
	return nil
}
