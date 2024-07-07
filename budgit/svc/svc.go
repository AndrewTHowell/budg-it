package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
)

type DB interface {
	InsertAccounts(ctx context.Context, accounts ...*budgit.Account) error
	SelectAccountByID(ctx context.Context, accountID string) (*budgit.Account, error)

	InsertExternalAccounts(ctx context.Context, externalAccounts ...*budgit.ExternalAccount) error
	SelectExternalAccountByID(ctx context.Context, externalAccountID string) (*budgit.ExternalAccount, error)
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

type Provider interface {
	ID() string
	GetExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error)
	GetExternalAccount(ctx context.Context, externalID string) (*budgit.ExternalAccount, error)
}

func (s Service) LoadAccountsFromProvider(ctx context.Context, budgetID, providerID string) ([]*budgit.Account, []*budgit.ExternalAccount, error) {
	externalAccounts, err := s.providers[providerID].GetExternalAccounts(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("loading accounts from %q: %w", providerID, err)
	}

	accounts := make([]*budgit.Account, 0, len(externalAccounts))
	for _, externalAccount := range externalAccounts {
		name := fmt.Sprintf("%s - %s", externalAccount.ExternalProviderID, externalAccount.Name)
		accounts = append(accounts, budgit.NewAccount(budgetID, name, externalAccount.ID, externalAccount.Balance))
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
	return fmt.Sprintf("Syncing Account %q failed, balance synced from external account %+v does not match balance of internal account %+v", e.AccountName, e.ExternalBalance, e.InternalBalance)
}

func (s Service) SyncAccount(ctx context.Context, budgetID, accountID string) error {
	account, err := s.db.SelectAccountByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("syncing account %q: %w", accountID, err)
	}
	currentExternalAccount, err := s.db.SelectExternalAccountByID(ctx, account.ExternalAccountID)
	if err != nil {
		return fmt.Errorf("syncing account %q: %w", accountID, err)
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
