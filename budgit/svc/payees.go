package svc

import (
	"context"
	"errors"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"golang.org/x/exp/maps"
)

type PayeeDB interface {
	InsertPayees(ctx context.Context, queryer db.Queryer, payee ...*budgit.Payee) error
	SelectPayeesByID(ctx context.Context, queryer db.Queryer, payeeIDs ...string) (map[string]*budgit.Payee, error)
	SelectPayeesByName(ctx context.Context, queryer db.Queryer, payeeNames ...string) (map[string]*budgit.Payee, error)
}

func (s Service) CreatePayees(ctx context.Context, payees ...*budgit.Payee) ([]*budgit.Payee, error) {
	if err := s.validatePayees(ctx, payees...); err != nil {
		return nil, fmt.Errorf("creating payees: %w", err)
	}

	if err := s.db.InsertPayees(ctx, payees...); err != nil {
		return nil, fmt.Errorf("creating payees: %w", err)
	}
	return payees, nil
}

type DuplicatePayeesError struct {
	PayeeNames []string
}

func (e DuplicatePayeesError) Error() string {
	return fmt.Sprintf("payees created with names that already exist: %+v", e.PayeeNames)
}

func (s Service) validatePayees(ctx context.Context, payees ...*budgit.Payee) error {
	errs := []error{}

	payeeNames := make([]string, 0, len(payees))
	for _, payee := range payees {
		payeeNames = append(payeeNames, payee.Name)
	}

	uniquePayeeNames := deduplicate(payeeNames)
	foundPayees, err := s.db.SelectPayeesByName(ctx, uniquePayeeNames...)
	if err != nil {
		return fmt.Errorf("validating payees: %w", err)
	}
	if len(foundPayees) < len(uniquePayeeNames) {
		foundPayeeNames := make([]string, 0, len(foundPayees))
		for _, payee := range foundPayees {
			foundPayeeNames = append(foundPayeeNames, payee.Name)
		}
		duplicateNames := intersection(uniquePayeeNames, maps.Keys(foundPayees))
		errs = append(errs, DuplicatePayeesError{PayeeNames: duplicateNames})
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}
