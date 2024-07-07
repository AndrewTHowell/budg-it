package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
)

type TransactionDB interface {
	InsertTransactions(ctx context.Context, transactions ...*budgit.Transaction) error
}

func (s Service) CreateTransactions(ctx context.Context, transactions ...*budgit.Transaction) ([]*budgit.Transaction, error) {

	if err := s.validateTransactions(ctx, transactions...); err != nil {
		return nil, fmt.Errorf("creating transactions: %w", err)
	}

	if err := s.db.InsertTransactions(ctx, transactions...); err != nil {
		return nil, fmt.Errorf("creating transactions: %w", err)
	}
	return transactions, nil
}

func (s Service) validateTransactions(_ context.Context, transactions ...*budgit.Transaction) error {
	for range transactions {
		break
	}
	return nil
}
