package db

import (
	"context"

	"github.com/andrewthowell/budgit/budgit"
)

type DB struct{}

func New() *DB {
	return &DB{}
}

func (db *DB) InsertAccounts(ctx context.Context, accounts ...*budgit.Account) error {
	return nil
}

func (db *DB) SelectAccountByID(ctx context.Context, accountID string) (*budgit.Account, error) {
	return nil, nil
}

func (db *DB) InsertExternalAccounts(ctx context.Context, accounts ...*budgit.ExternalAccount) error {
	return nil
}

func (db *DB) SelectExternalAccountByID(ctx context.Context, externalAccountID string) (*budgit.ExternalAccount, error) {
	return nil, nil
}
