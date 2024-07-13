package svc

import (
	"context"
	"fmt"
	"time"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db/dbconvert"
	"github.com/google/uuid"
)

type Integration interface {
	ID() string
	GetExternalAccounts(ctx context.Context, syncTime time.Time) ([]*budgit.ExternalAccount, error)
	GetExternalAccount(ctx context.Context, syncTime time.Time, externalID string) (*budgit.ExternalAccount, error)
}

func (s Service) LoadAccountsFromIntegration(ctx context.Context, integrationID string) ([]*budgit.Account, error) {
	externalAccounts, err := s.integrations[integrationID].GetExternalAccounts(ctx, time.Now().UTC())
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

	// TODO: check for accounts not being inserted
	if _, err := s.db.InsertAccounts(ctx, s.conn, dbconvert.FromAccounts(accounts...)...); err != nil {
		return nil, fmt.Errorf("loading accounts from %q: %w", integrationID, err)
	}
	return accounts, nil
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

	externalAccount, err := s.integrations[account.ExternalAccount.IntegrationID].GetExternalAccount(ctx, time.Now().UTC(), account.ExternalAccount.ID)
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
	account.ExternalAccount.LastSyncTimestamp = time.Now().UTC()
	account.ExternalAccount = externalAccount

	// TODO: check for accounts not being inserted
	if _, err := s.db.InsertAccounts(ctx, s.conn, dbconvert.FromAccounts(account)...); err != nil {
		return fmt.Errorf("syncing account %q: %w", accountID, err)
	}
	return nil
}
