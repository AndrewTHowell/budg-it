package svc

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/db/dbconvert"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/maps"
)

type TransactionDB interface {
	InsertTransactions(ctx context.Context, queryer db.Queryer, transactions ...*db.Transaction) ([]string, error)
}

func (s Service) CreateTransactions(ctx context.Context, transactions ...*budgit.Transaction) ([]*budgit.Transaction, error) {
	var createdTransactions []*budgit.Transaction
	err := s.inTx(ctx, func(conn Conn) error {
		if err := s.validateTransactions(ctx, conn, transactions...); err != nil {
			return err
		}

		transactions, err := s.appendMirrorTransactions(transactions...)
		if err != nil {
			return err
		}

		// TODO: check for transactions not being inserted
		if _, err := s.db.InsertTransactions(ctx, conn, dbconvert.FromTransactions(transactions...)...); err != nil {
			return err
		}
		createdTransactions = transactions
		// Update Account/Category Balances
		return nil
	}, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		return nil, fmt.Errorf("creating transactions: %w", err)
	}
	return createdTransactions, nil
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

func (s Service) validateTransactions(ctx context.Context, conn Conn, transactions ...*budgit.Transaction) error {
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
	foundAccounts, err := s.db.SelectAccountsByID(ctx, conn, uniqueAccountIDs...)
	if err != nil {
		return fmt.Errorf("validating transactions: %w", err)
	}
	if len(foundAccounts) < len(uniqueAccountIDs) {
		missingIDs := symmetricDifference(uniqueAccountIDs, maps.Keys(foundAccounts))
		errs = append(errs, MissingAccountsError{AccountIDs: missingIDs})
	}

	uniquePayeeIDs := deduplicate(payeeIDs)
	foundPayees, err := s.db.SelectPayeesByID(ctx, conn, uniquePayeeIDs...)
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

func (s Service) appendMirrorTransactions(transactions ...*budgit.Transaction) ([]*budgit.Transaction, error) {
	mirrorTransactions := make([]*budgit.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		if transaction.IsPayeeInternal {
			mirrorTransactions = append(mirrorTransactions, transaction.Mirror(uuid.New().String()))
		}
	}
	return slices.Concat(transactions, mirrorTransactions), nil
}
