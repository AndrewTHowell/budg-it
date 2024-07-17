package svc

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type Conn interface {
	db.Execer
	db.Queryer
}

type TxConn interface {
	Conn
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type DB interface {
	Now(ctx context.Context, queryer db.Queryer) (pgtype.Timestamptz, error)
	AccountDB
	PayeeDB
	TransactionDB
}

type Service struct {
	log          *zap.SugaredLogger
	conn         TxConn
	db           DB
	integrations map[string]Integration
}

func New(log *zap.SugaredLogger, conn TxConn, db DB, integrations []Integration) *Service {
	return &Service{
		log:          log,
		conn:         conn,
		db:           db,
		integrations: mapByID(integrations),
	}
}

func (s Service) inTx(ctx context.Context, txFunc func(conn Conn) error, txOptions pgx.TxOptions) (rollbackErr error) {
	// rollbackErr is a named return so that it can be modified in a deferred call.

	tx, err := s.conn.BeginTx(ctx, txOptions)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			rollbackErr = fmt.Errorf("failed to rollback transaction: %w", err)
		}
	}()

	if err := txFunc(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	// Reaching here means all normal errors have been avoided, but must return rollbackErr in case defer errors and the error must be returned.
	return rollbackErr
}

type idGetter interface {
	ID() string
}

func mapByID[E idGetter](elems []E) map[string]E {
	elemsByID := make(map[string]E, len(elems))
	for _, elem := range elems {
		elemsByID[elem.ID()] = elem
	}
	return elemsByID
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
