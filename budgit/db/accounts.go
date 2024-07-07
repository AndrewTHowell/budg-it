package db

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
)

func (db *DB) InsertAccounts(ctx context.Context, accounts ...*budgit.Account) error {
	if err := insert(&db.accounts, accounts...); err != nil {
		return fmt.Errorf("inserting %d accounts: %w", len(accounts), err)
	}
	return nil
}

func (db *DB) SelectAccountsByID(ctx context.Context, accountIDs ...string) (map[string]*budgit.Account, error) {
	return selectByIDs(db.accounts, accountIDs), nil
}

func (db *DB) SelectAccounts(ctx context.Context) ([]*budgit.Account, error) {
	return db.accounts, nil
}

func (db *DB) InsertExternalAccounts(ctx context.Context, externalAccounts ...*budgit.ExternalAccount) error {
	if err := insert(&db.externalAccounts, externalAccounts...); err != nil {
		return fmt.Errorf("inserting %d external accounts: %w", len(externalAccounts), err)
	}
	return nil
}

func (db *DB) SelectExternalAccountsByID(ctx context.Context, externalAccountIDs ...string) (map[string]*budgit.ExternalAccount, error) {
	return selectByIDs(db.externalAccounts, externalAccountIDs), nil
}

func (db *DB) SelectExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error) {
	return db.externalAccounts, nil
}
