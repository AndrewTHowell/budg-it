package db

import (
	"context"

	"github.com/andrewthowell/budgit/budgit"
)

type DB struct{}

func New() *DB {
	return &DB{}
}

func (db *DB) InsertAccounts(ctx context.Context, accounts []*budgit.Account) error {
	return nil
}

func (db *DB) InsertExternalAccounts(ctx context.Context, accounts []*budgit.ExternalAccount) error {
	return nil
}
