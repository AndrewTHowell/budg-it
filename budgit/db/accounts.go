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

var ErrAccountNotFound = fmt.Errorf("the requested Account does not exist")

func (db *DB) SelectAccountByID(ctx context.Context, accountID string) (*budgit.Account, error) {
	account, err := selectByID(db.accounts, accountID, ErrAccountNotFound)
	if err != nil {
		return nil, fmt.Errorf("selecting account by ID %q: %w", accountID, err)
	}
	return account, nil
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

var ErrExternalAccountNotFound = fmt.Errorf("the requested ExternalAccount does not exist")

func (db *DB) SelectExternalAccountByID(ctx context.Context, externalAccountID string) (*budgit.ExternalAccount, error) {
	externalAccount, err := selectByID(db.externalAccounts, externalAccountID, ErrExternalAccountNotFound)
	if err != nil {
		return nil, fmt.Errorf("selecting external account by ID %q: %w", externalAccountID, err)
	}
	return externalAccount, nil
}

func (db *DB) SelectExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error) {
	return db.externalAccounts, nil
}
