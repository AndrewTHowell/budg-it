package db

import (
	"context"
	"fmt"
	"slices"

	"github.com/andrewthowell/budgit/budgit"
)

type DB struct {
	accounts         []*budgit.Account
	externalAccounts []*budgit.ExternalAccount
}

func New() *DB {
	return &DB{
		accounts:         []*budgit.Account{},
		externalAccounts: []*budgit.ExternalAccount{},
	}
}

func (db *DB) InsertAccounts(ctx context.Context, accounts ...*budgit.Account) error {
	if err := insert(db.accounts, accounts...); err != nil {
		return fmt.Errorf("inserting %d accounts: %w", len(accounts), err)
	}
	return nil
}

var ErrAccountNotFound = fmt.Errorf("the requested Account does not exist")

func (db *DB) SelectAccountByID(ctx context.Context, accountID string) (*budgit.Account, error) {
	account, err := selectByID(db.accounts, accountID, ErrAccountNotFound)
	if err != nil {
		return nil, fmt.Errorf("selecting account by ID %q: %w", accountID, err)
	}
	return account, nil
}

func (db *DB) InsertExternalAccounts(ctx context.Context, externalAccounts ...*budgit.ExternalAccount) error {
	if err := insert(db.externalAccounts, externalAccounts...); err != nil {
		return fmt.Errorf("inserting %d external accounts: %w", len(externalAccounts), err)
	}
	return nil
}

var ErrExternalAccountNotFound = fmt.Errorf("the requested ExternalAccount does not exist")

func (db *DB) SelectExternalAccountByID(ctx context.Context, externalAccountID string) (*budgit.ExternalAccount, error) {
	externalAccount, err := selectByID(db.externalAccounts, externalAccountID, ErrExternalAccountNotFound)
	if err != nil {
		return nil, fmt.Errorf("selecting external account by ID %q: %w", externalAccountID, err)
	}
	return externalAccount, nil
}

type idGetter interface {
	GetID() string
}

func insert[I idGetter](slice []*I, elems ...*I) error {
	slice = slices.Grow(slice, len(elems))
	for _, elem := range elems {
		idx := slices.IndexFunc(slice, func(i *I) bool {
			return (*elem).GetID() < (*i).GetID()
		})
		if idx != -1 {
			slice = slices.Insert(slice, idx, elem)
		}
		// If no index can be found where given element has lesser ID, element has the greatest ID and belongs at the end.
		slice = append(slice, elem)
	}
	return nil
}

func selectByID[I idGetter](slice []*I, id string, notFoundError error) (*I, error) {
	idx, ok := slices.BinarySearchFunc(slice, id, func(elem *I, targetID string) int {
		if (*elem).GetID() < targetID {
			return -1
		}
		if (*elem).GetID() == targetID {
			return 0
		}
		return 1
	})
	if !ok {
		return nil, notFoundError
	}
	return slice[idx], nil
}
