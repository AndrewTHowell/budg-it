package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/db/dbconvert"
)

type AccountDB interface {
	InsertAccounts(ctx context.Context, queryer db.Queryer, account ...*db.Account) ([]string, error)
	SelectAccounts(ctx context.Context, queryer db.Queryer) ([]*db.Account, error)
	SelectAccountsByID(ctx context.Context, queryer db.Queryer, accountIDs ...string) (map[string]*db.Account, error)
}

func (s Service) ListAccounts(ctx context.Context) ([]*budgit.Account, error) {
	accounts, err := s.db.SelectAccounts(ctx, s.conn)
	if err != nil {
		return nil, fmt.Errorf("listing accounts: %w", err)
	}
	return dbconvert.ToAccounts(accounts...), nil
}
