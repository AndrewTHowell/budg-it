package db

import (
	"slices"

	"github.com/andrewthowell/budgit/budgit"
)

type DB struct {
	accounts         []*budgit.Account
	externalAccounts []*budgit.ExternalAccount
	transactions     []*budgit.Transaction
}

func New() *DB {
	return &DB{
		accounts:         []*budgit.Account{},
		externalAccounts: []*budgit.ExternalAccount{},
		transactions:     []*budgit.Transaction{},
	}
}

type idGetter interface {
	GetID() string
}

func insert[I idGetter](slice *[]*I, elems ...*I) error {
	*slice = slices.Grow(*slice, len(elems))
	for _, elem := range elems {
		idx := slices.IndexFunc(*slice, func(i *I) bool {
			return (*elem).GetID() < (*i).GetID()
		})
		if idx != -1 {
			*slice = slices.Insert(*slice, idx, elem)
		}
		// If no index can be found where given element has lesser ID, element has the greatest ID and belongs at the end.
		*slice = append(*slice, elem)
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
