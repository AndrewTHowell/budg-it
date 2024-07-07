package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
)

func (s Service) ListAccounts(ctx context.Context) ([]*budgit.Account, error) {
	accounts, err := s.db.SelectAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing accounts: %w", err)
	}
	return accounts, nil
}

func (s Service) ListExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error) {
	externalAccounts, err := s.db.SelectExternalAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing external accounts: %w", err)
	}
	return externalAccounts, nil
}
