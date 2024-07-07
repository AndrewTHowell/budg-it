package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
)

type Provider interface {
	ID() string
	GetExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error)
	GetExternalAccount(ctx context.Context, externalID string) (*budgit.ExternalAccount, error)
}

func (s Service) LoadAccountsFromProvider(ctx context.Context, providerID string) ([]*budgit.Account, []*budgit.ExternalAccount, error) {
	externalAccounts, err := s.providers[providerID].GetExternalAccounts(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("loading accounts from %q: %w", providerID, err)
	}

	accounts := make([]*budgit.Account, 0, len(externalAccounts))
	for _, externalAccount := range externalAccounts {
		name := fmt.Sprintf("%s - %s", externalAccount.ExternalProviderID, externalAccount.Name)
		accounts = append(accounts, budgit.NewAccount(name, externalAccount.ID, externalAccount.Balance))
	}

	if err := s.db.InsertAccounts(ctx, accounts...); err != nil {
		return nil, nil, fmt.Errorf("loading accounts from %q: %w", providerID, err)
	}
	if err := s.db.InsertExternalAccounts(ctx, externalAccounts...); err != nil {
		return nil, nil, fmt.Errorf("loading accounts from %q: %w", providerID, err)
	}
	return accounts, externalAccounts, nil
}

type AccountSyncError struct {
	AccountName                      string
	ExternalBalance, InternalBalance budgit.Balance
}

func (e AccountSyncError) Error() string {
	return fmt.Sprintf("syncing Account %q failed, balance synced from external account %+v does not match balance of internal account %+v", e.AccountName, e.ExternalBalance, e.InternalBalance)
}

var ErrAccountNotFound = fmt.Errorf("the requested Account does not exist")

func (s Service) SyncAccount(ctx context.Context, accountID string) error {
	accounts, err := s.db.SelectAccountsByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("syncing account %q: %w", accountID, err)
	}
	account, ok := accounts[accountID]
	if !ok {
		return fmt.Errorf("syncing account %q: %w", accountID, ErrAccountNotFound)
	}
	currentExternalAccounts, err := s.db.SelectExternalAccountsByID(ctx, account.ExternalAccountID)
	if err != nil {
		return fmt.Errorf("syncing account %q: %w", accountID, err)
	}
	currentExternalAccount, ok := currentExternalAccounts[account.ExternalAccountID]
	if !ok {
		return fmt.Errorf("syncing account %q: account references external account %q which does not exist", accountID, account.ExternalAccountID)
	}
	externalAccount, err := s.providers[currentExternalAccount.ExternalProviderID].GetExternalAccount(ctx, account.ExternalAccountID)
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
	if err := s.db.InsertExternalAccounts(ctx, externalAccount); err != nil {
		return fmt.Errorf("syncing account %q: %w", accountID, err)
	}
	return nil
}
