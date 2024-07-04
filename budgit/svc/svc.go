package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
)

type DB interface {
	InsertAccounts(ctx context.Context, accounts []*budgit.Account) error
	InsertExternalAccounts(ctx context.Context, externalAccounts []*budgit.ExternalAccount) error
}

type Service struct {
	db DB
}

func New(db DB) *Service {
	return &Service{
		db: db,
	}
}

type Provider interface {
	ID() string
	GetExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error)
	GetExternalAccount(ctx context.Context, externalID string) (*budgit.ExternalAccount, error)
}

func (s Service) LoadAccountsFromProvider(ctx context.Context, provider Provider, budgetID string) error {
	externalAccounts, err := provider.GetExternalAccounts(ctx)
	if err != nil {
		return fmt.Errorf("loading accounts from %s: %w", provider.ID(), err)
	}

	accounts := make([]*budgit.Account, 0, len(externalAccounts))
	for _, externalAccount := range externalAccounts {
		name := fmt.Sprintf("%s - %s", externalAccount.ExternalProviderID, externalAccount.Name)
		accounts = append(accounts, budgit.NewAccount(budgetID, name, externalAccount.ID))
	}

	if err := s.db.InsertAccounts(ctx, accounts); err != nil {
		return fmt.Errorf("loading accounts from %s: %w", provider.ID(), err)
	}
	if err := s.db.InsertExternalAccounts(ctx, externalAccounts); err != nil {
		return fmt.Errorf("loading accounts from %s: %w", provider.ID(), err)
	}
	return nil
}
