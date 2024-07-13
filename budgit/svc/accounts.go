package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
)

type AccountDB interface {
	InsertAccounts(ctx context.Context, queryer db.Queryer, account ...*db.Account) error
	SelectAccounts(ctx context.Context, queryer db.Queryer) ([]*db.Account, error)
	SelectAccountsByID(ctx context.Context, queryer db.Queryer, accountIDs ...string) (map[string]*db.Account, error)
	SelectAccountsByName(ctx context.Context, queryer db.Queryer, accountNames ...string) (map[string]*d.Account, error)
}

func (s Service) ListAccounts(ctx context.Context) ([]*budgit.Account, error) {
	accounts, err := s.db.SelectAccounts(ctx, s.conn)
	if err != nil {
		return nil, fmt.Errorf("listing accounts: %w", err)
	}
	return accounts, nil
}
