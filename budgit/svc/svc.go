package svc

import (
	"context"

	"github.com/andrewthowell/budgit/budgit"
)

type DB interface {
	InsertAccounts(ctx context.Context, accounts ...*budgit.Account) error
	SelectAccountsByID(ctx context.Context, accountIDs ...string) (map[string]*budgit.Account, error)
	SelectAccounts(ctx context.Context) ([]*budgit.Account, error)

	InsertExternalAccounts(ctx context.Context, externalAccounts ...*budgit.ExternalAccount) error
	SelectExternalAccountsByID(ctx context.Context, externalAccountIDs ...string) (map[string]*budgit.ExternalAccount, error)
	SelectExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error)

	TransactionDB
	PayeeDB
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

func deduplicate[S ~[]E, E comparable](slice S) S {
	seen := make(map[E]bool, len(slice))
	deduplicated := make([]E, 0, len(slice))
	for _, e := range slice {
		if !seen[e] {
			seen[e] = true
			deduplicated = append(deduplicated, e)
		}
	}
	return deduplicated
}

func symmetricDifference[S ~[]E, E comparable](sliceA, sliceB S) S {
	seenA := boolMap(sliceA)
	diff := make([]E, 0, len(sliceA)+len(sliceB))
	for _, b := range sliceB {
		if !seenA[b] {
			diff = append(diff, b)
		}
	}
	return diff
}

func intersection[S ~[]E, E comparable](sliceA, sliceB S) S {
	seenA := boolMap(sliceA)
	intersect := make([]E, 0, len(sliceA)+len(sliceB))
	for _, b := range sliceB {
		if seenA[b] {
			intersect = append(intersect, b)
		}
	}
	return intersect
}

func boolMap[E comparable](slice []E) map[E]bool {
	m := make(map[E]bool, len(slice))
	for _, e := range slice {
		m[e] = true
	}
	return m
}
