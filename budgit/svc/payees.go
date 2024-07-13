package svc

import (
	"context"
	"errors"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/db/dbconvert"
	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/maps"
)

type PayeeDB interface {
	InsertPayees(ctx context.Context, queryer db.Queryer, payee ...*db.Payee) ([]string, error)
	SelectPayeesByID(ctx context.Context, queryer db.Queryer, payeeIDs ...string) (map[string]*db.Payee, error)
	SelectPayeesByName(ctx context.Context, queryer db.Queryer, payeeNames ...string) (map[string]*db.Payee, error)
}

func (s Service) CreatePayees(ctx context.Context, payees ...*budgit.Payee) ([]*budgit.Payee, error) {
	var createdPayees []*budgit.Payee
	err := s.inTx(ctx, func(conn Conn) error {
		if err := s.validatePayees(ctx, conn, payees...); err != nil {
			return err
		}

		// TODO: check for payees not being inserted
		if _, err := s.db.InsertPayees(ctx, conn, dbconvert.FromPayees(payees...)...); err != nil {
			return err
		}
		createdPayees = payees
		return nil
	}, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		return nil, fmt.Errorf("creating payees: %w", err)
	}
	return createdPayees, nil
}

type DuplicatePayeesError struct {
	PayeeNames []string
}

func (e DuplicatePayeesError) Error() string {
	return fmt.Sprintf("payees created with names that already exist: %+v", e.PayeeNames)
}

func (s Service) validatePayees(ctx context.Context, conn Conn, payees ...*budgit.Payee) error {
	errs := []error{}

	payeeNames := make([]string, 0, len(payees))
	for _, payee := range payees {
		payeeNames = append(payeeNames, payee.Name)
	}

	uniquePayeeNames := deduplicate(payeeNames)
	foundDBPayees, err := s.db.SelectPayeesByName(ctx, conn, uniquePayeeNames...)
	if err != nil {
		return fmt.Errorf("validating payees: %w", err)
	}
	if len(foundDBPayees) < len(uniquePayeeNames) {
		foundPayeeNames := make([]string, 0, len(foundDBPayees))
		for _, dbPayee := range foundDBPayees {
			foundPayeeNames = append(foundPayeeNames, dbconvert.ToPayees(dbPayee)[0].Name)
		}
		duplicateNames := intersection(uniquePayeeNames, maps.Keys(foundDBPayees))
		errs = append(errs, DuplicatePayeesError{PayeeNames: duplicateNames})
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}
