package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Payee struct {
	ID   pgtype.Text `db:"id"`
	Name pgtype.Text `db:"name"`
}

func (p Payee) GetID() string {
	return p.ID.String
}

func (p Payee) GetName() string {
	return p.Name.String
}

var (
	payeeColumns    = getAllDBColumns(Payee{})
	payeeColumnsStr = strings.Join(payeeColumns, ", ")
)

func (db DB) InsertPayees(ctx context.Context, queryer Queryer, payees ...*Payee) ([]string, error) {
	sql := fmt.Sprintf(`
		INSERT INTO payees (%[1]s)
		(
			SELECT %[1]s
			FROM UNNEST(
				$1::TEXT[],
				$2::TEXT[]
			)
			AS u(%[1]s)
		)
		ON CONFLICT DO NOTHING
		RETURNING id;
	`, payeeColumnsStr)

	rows, err := queryer.Query(ctx, sql, payeesToArgs(payees)...)
	if err != nil {
		return nil, fmt.Errorf("inserting %d payees: %w", len(payees), err)
	}
	defer rows.Close()

	ids, err := rowsToIDs(rows)
	if err != nil {
		return nil, fmt.Errorf("inserting %d payees: %w", len(payees), err)
	}
	return ids, nil
}

func (db DB) SelectPayees(ctx context.Context, queryer Queryer) ([]*Payee, error) {
	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM payees
		ORDER BY id
	`, payeeColumnsStr)

	rows, err := queryer.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("selecting payees: %w", err)
	}
	defer rows.Close()

	payees, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payee])
	if err != nil {
		return nil, fmt.Errorf("selecting payees: %w", err)
	}
	return structsToPointers(payees), nil
}

func (db DB) SelectPayeesByID(ctx context.Context, queryer Queryer, payeeIDs ...string) (map[string]*Payee, error) {
	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM payees
		WHERE id = ANY($1::TEXT[])
	`, payeeColumnsStr)

	ids := make([]pgtype.Text, 0, len(payeeIDs))
	for _, id := range payeeIDs {
		ids = append(ids, pgtype.Text{String: id, Valid: true})
	}

	rows, err := queryer.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("selecting payees by ID: %w", err)
	}
	defer rows.Close()

	payees, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payee])
	if err != nil {
		return nil, fmt.Errorf("selecting payees by ID: %w", err)
	}
	return mapByID(structsToPointers(payees)), nil
}

func (db DB) SelectPayeesByName(ctx context.Context, queryer Queryer, payeeNames ...string) (map[string]*Payee, error) {
	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM payees
		WHERE name = ANY($1::TEXT[])
	`, payeeColumnsStr)

	names := make([]pgtype.Text, 0, len(payeeNames))
	for _, name := range payeeNames {
		names = append(names, pgtype.Text{String: name, Valid: true})
	}

	rows, err := queryer.Query(ctx, sql, names)
	if err != nil {
		return nil, fmt.Errorf("selecting payees by name: %w", err)
	}
	defer rows.Close()

	payees, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payee])
	if err != nil {
		return nil, fmt.Errorf("selecting payees by name: %w", err)
	}
	return mapByName(structsToPointers(payees)), nil
}

func payeesToArgs(payees []*Payee) []any {
	ids := make([]pgtype.Text, 0, len(payees))
	names := make([]pgtype.Text, 0, len(payees))
	for _, payee := range payees {
		ids = append(ids, payee.ID)
		names = append(names, payee.Name)
	}
	return []any{
		ids,
		names,
	}
}
