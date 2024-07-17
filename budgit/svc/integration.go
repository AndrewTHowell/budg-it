package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/db/dbconvert"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Integration interface {
	ID() string
	GetExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error)
	GetExternalAccount(ctx context.Context, externalID string) (*budgit.ExternalAccount, error)
}

func (s Service) LoadAccountsFromIntegration(ctx context.Context, integrationID string) ([]*budgit.Account, error) {
	externalAccounts, err := s.integrations[integrationID].GetExternalAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading accounts from %q: %w", integrationID, err)
	}

	accounts := make([]*budgit.Account, 0, len(externalAccounts))
	for _, externalAccount := range externalAccounts {
		accounts = append(accounts, &budgit.Account{
			ID:              uuid.New().String(),
			Name:            fmt.Sprintf("%s - %s", externalAccount.IntegrationID, externalAccount.Name),
			Balance:         externalAccount.Balance,
			ExternalAccount: externalAccount,
		})
	}

	var createdAccounts []*budgit.Account
	err = s.inTx(ctx, func(conn Conn) error {
		now, err := s.db.Now(ctx, conn)
		if err != nil {
			return err
		}

		dbAccounts := dbconvert.FromAccounts(accounts...)
		for _, dbAccount := range dbAccounts {
			dbAccount.ValidFromTimestamp = now
			dbAccount.ValidToTimestamp = pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true}
			dbAccount.ExternalLastSyncTimestamp = now
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

type AccountSyncError struct {
	AccountName                      string
	ExternalBalance, InternalBalance budgit.Balance
}

func (e AccountSyncError) Error() string {
	return fmt.Sprintf("syncing Account %q failed, balance synced from external account %+v does not match balance of internal account %+v", e.AccountName, e.ExternalBalance, e.InternalBalance)
}

var (
	ErrAccountNotFound  = fmt.Errorf("the requested Account does not exist")
	ErrAccountNotLinked = fmt.Errorf("the requested Account is not linked to an external account and cannot be synced")
)

func (s Service) SyncAccount(ctx context.Context, accountID string) error {
	dbAccounts, err := s.db.SelectAccountsByID(ctx, s.conn, accountID)
	if err != nil {
		return fmt.Errorf("syncing account %q: %w", accountID, err)
	}
	dbAccount, ok := dbAccounts[accountID]
	if !ok {
		return fmt.Errorf("syncing account %q: %w", accountID, ErrAccountNotFound)
	}
	account := dbconvert.ToAccounts(dbAccount)[0]

	if account.ExternalAccount == nil {
		// Account is not linked, no need to sync
		return ErrAccountNotLinked
	}

	externalAccount, err := s.integrations[account.ExternalAccount.IntegrationID].GetExternalAccount(ctx, account.ExternalAccount.ID)
	if err != nil {
		return fmt.Errorf("syncing account %q: %w", accountID, err)
	}
	if account.Balance != externalAccount.Balance {
		return AccountSyncError{
			AccountName:     account.Name,
			ExternalBalance: externalAccount.Balance,
			InternalBalance: account.Balance,
		}
	}
	account.ExternalAccount = externalAccount

	err = s.inTx(ctx, func(conn Conn) error {
		now, err := s.db.Now(ctx, conn)
		if err != nil {
			return err
		}

		if _, err := s.db.UpdateAccountValidToTimestamps(ctx, s.conn, db.ValidToTimestampUpdate{
			ID:               dbAccount.ID,
			ValidToTimestamp: now,
		}); err != nil {
			return fmt.Errorf("syncing account %q: %w", accountID, err)
		}

		dbAccount := dbconvert.FromAccounts(account)[0]
		dbAccount.ValidFromTimestamp = now
		dbAccount.ValidToTimestamp = pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true}
		dbAccount.ExternalLastSyncTimestamp = now

		// TODO: check for accounts not being inserted
		if _, err := s.db.InsertAccounts(ctx, s.conn, dbAccount); err != nil {
			return fmt.Errorf("syncing account %q: %w", accountID, err)
		}
		return nil
	}, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		return fmt.Errorf("syncing account %q: %w", accountID, err)
	}
	return nil
}
