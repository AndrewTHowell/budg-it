package svc

import (
	"context"

	"github.com/andrewthowell/budgit/budgit"
)

type DB interface {
	InsertAccounts(ctx context.Context, accounts ...*budgit.Account) error
	SelectAccountByID(ctx context.Context, accountID string) (*budgit.Account, error)
	SelectAccounts(ctx context.Context) ([]*budgit.Account, error)

	InsertExternalAccounts(ctx context.Context, externalAccounts ...*budgit.ExternalAccount) error
	SelectExternalAccountByID(ctx context.Context, externalAccountID string) (*budgit.ExternalAccount, error)
	SelectExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error)

	TransactionDB
}

type Service struct {
	db        DB
	providers map[string]Provider
}

func New(db DB, providers map[string]Provider) *Service {
	return &Service{
		db:        db,
		providers: providers,
	}
}
