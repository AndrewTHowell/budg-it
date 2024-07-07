package svc

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/andrewthowell/budgit/budgit"
	"golang.org/x/exp/maps"
)

type TransactionDB interface {
	InsertTransactions(ctx context.Context, transactions ...*budgit.Transaction) error
}

func (s Service) CreateTransactions(ctx context.Context, transactions ...*budgit.Transaction) ([]*budgit.Transaction, error) {
	if err := s.validateTransactions(ctx, transactions...); err != nil {
		return nil, fmt.Errorf("creating transactions: %w", err)
	}

	transactions, err := s.addMirrorTransactions(ctx, transactions...)
	if err != nil {
		return nil, fmt.Errorf("creating transactions: %w", err)
	}

	if err := s.db.InsertTransactions(ctx, transactions...); err != nil {
		return nil, fmt.Errorf("creating transactions: %w", err)
	}
	// Update Account/Category Balances

	return transactions, nil
}

type MissingAccountsError struct {
	AccountIDs []string
}

func (e MissingAccountsError) Error() string {
	return fmt.Sprintf("transactions reference Accounts that do not exist: %+v", e.AccountIDs)
}

type MissingPayeesError struct {
	PayeeIDs []string
}

func (e MissingPayeesError) Error() string {
	return fmt.Sprintf("transactions reference Payees that do not exist: %+v", e.PayeeIDs)
}

func (s Service) validateTransactions(ctx context.Context, transactions ...*budgit.Transaction) error {
	errs := []error{}

	accountIDs := make([]string, 0, len(transactions))
	payeeIDs := make([]string, 0, len(transactions))
	for _, transaction := range transactions {
		accountIDs = append(accountIDs, transaction.AccountID)
		if transaction.IsPayeeInternal {
			accountIDs = append(accountIDs, transaction.PayeeID)
		} else {
			payeeIDs = append(payeeIDs, transaction.PayeeID)
		}
	}

	uniqueAccountIDs := deduplicate(accountIDs)
	foundAccounts, err := s.db.SelectAccountsByID(ctx, uniqueAccountIDs...)
	if err != nil {
		return fmt.Errorf("validating transactions: %w", err)
	}
	if len(foundAccounts) < len(uniqueAccountIDs) {
		missingIDs := symmetricDifference(uniqueAccountIDs, maps.Keys(foundAccounts))
		errs = append(errs, MissingAccountsError{AccountIDs: missingIDs})
	}

	uniquePayeeIDs := deduplicate(payeeIDs)
	foundPayees, err := s.db.SelectPayeesByID(ctx, uniquePayeeIDs...)
	if err != nil {
		return fmt.Errorf("validating transactions: %w", err)
	}
	if len(foundPayees) < len(uniquePayeeIDs) {
		missingIDs := symmetricDifference(uniquePayeeIDs, maps.Keys(foundPayees))
		errs = append(errs, MissingPayeesError{PayeeIDs: missingIDs})
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (s Service) addMirrorTransactions(_ context.Context, transactions ...*budgit.Transaction) ([]*budgit.Transaction, error) {
	mirrorTransactions := make([]*budgit.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		if transaction.IsPayeeInternal {
			mirrorTransactions = append(mirrorTransactions, transaction.Mirror())
		}
	}
	return slices.Concat(transactions, mirrorTransactions), nil
}
