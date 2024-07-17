package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/db/dbconvert"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccountDB interface {
	InsertAccounts(ctx context.Context, queryer db.Queryer, account ...*db.Account) ([]string, error)
	SelectAccounts(ctx context.Context, queryer db.Queryer) ([]*db.Account, error)
	SelectAccountsByID(ctx context.Context, queryer db.Queryer, accountIDs ...string) (map[string]*db.Account, error)
}

func (s Service) CreateAccounts(ctx context.Context, accounts ...*budgit.Account) ([]*budgit.Account, error) {

	var createdAccounts []*budgit.Account
	err := s.inTx(ctx, func(conn Conn) error {
		now, err := s.db.Now(ctx, conn)
		if err != nil {
			return err
		}

		dbAccounts := dbconvert.FromAccounts(accounts...)
		for _, dbAccount := range dbAccounts {
			dbAccount.ValidFromTimestamp = now
			dbAccount.ValidToTimestamp = pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true}
		}

		// TODO: check for accounts not being inserted
		if _, err := s.db.InsertAccounts(ctx, conn, dbAccounts...); err != nil {
			return err
		}
		createdAccounts = accounts
		return nil
	}, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		return nil, fmt.Errorf("creating accounts: %w", err)
	}
	return createdAccounts, nil
}

func (s Service) ListAccounts(ctx context.Context) ([]*budgit.Account, error) {
	accounts, err := s.db.SelectAccounts(ctx, s.conn)
	if err != nil {
		return nil, fmt.Errorf("listing accounts: %w", err)
	}
	return dbconvert.ToAccounts(accounts...), nil
}
