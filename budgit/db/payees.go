package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type Payee struct {
	RequestID          pgtype.Text        `db:"request_id"`
	ValidFromTimestamp pgtype.Timestamptz `db:"valid_from_timestamp"`
	ValidToTimestamp   pgtype.Timestamptz `db:"valid_to_timestamp"`
	ID                 pgtype.Text        `db:"id"`
	Name               pgtype.Text        `db:"name"`
}

func (p Payee) GetID() string {
	return p.ID.String
}

func (p Payee) GetRequestID() string {
	return p.RequestID.String
}

func (p Payee) GetName() string {
	return p.Name.String
}

var (
	payeeColumns    = getAllDBColumns(Payee{})
	payeeColumnsStr = strings.Join(payeeColumns, ", ")
)

func (db DB) InsertPayees(ctx context.Context, queryer Queryer, payees ...*Payee) ([]string, error) {
	db.log.Debugw("Inserting payees", zap.Int("number_of_payees", len(payees)))

	sql := fmt.Sprintf(`
		INSERT INTO payees (%[1]s)
		(
			SELECT %[1]s
			FROM UNNEST(
				$1::TEXT[],
				$2::TIMESTAMPTZ[],
				$3::TIMESTAMPTZ[],
				$4::TEXT[],
				$5::TEXT[]
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
	db.log.Debugw("Inserted payees", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	ids, err := rowsToIDs(rows)
	if err != nil {
		return nil, fmt.Errorf("inserting %d payees: %w", len(payees), err)
	}
	db.log.Debugw("Inserted payees scanned", zap.String("inserted_ids", fmt.Sprintf("%v", ids)))
	return ids, nil
}

func (db DB) UpdatePayeeValidToTimestamps(ctx context.Context, queryer Queryer, updates ...ValidToTimestampUpdate) ([]string, error) {
	db.log.Debugw("Updating payee valid to timestamps", zap.Int("number_of_payees", len(updates)))

	sql := `
		UPDATE payees
		SET valid_to_timestamp = input.valid_to_timestamp
		FROM 
		(
			SELECT id, valid_to_timestamp
			FROM UNNEST(
				$1::TEXT[],
				$2::TIMESTAMPTZ[]
			)
			AS u(id, valid_to_timestamp)
		) AS input
		WHERE payees.valid_to_timestamp = 'infinity'
		AND payees.id = input.id
		RETURNING payees.id;
	`

	payeeIDs := make([]pgtype.Text, 0, len(updates))
	validToTimestamps := make([]pgtype.Timestamptz, 0, len(updates))
	for _, update := range updates {
		payeeIDs = append(payeeIDs, update.ID)
		validToTimestamps = append(validToTimestamps, update.ValidToTimestamp)
	}

	rows, err := queryer.Query(ctx, sql, payeeIDs, validToTimestamps)
	if err != nil {
		return nil, fmt.Errorf("updating %d payee valid to timestamps: %w", len(updates), err)
	}
	defer rows.Close()
	db.log.Debugw("Updated payee valid to timestamps", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	ids, err := rowsToIDs(rows)
	if err != nil {
		return nil, fmt.Errorf("updating %d payee valid to timestamps: %w", len(updates), err)
	}
	db.log.Debugw("Updated payee valid to timestamps scanned", zap.String("updated_ids", fmt.Sprintf("%v", ids)))
	return ids, nil
}

func (db DB) SelectPayees(ctx context.Context, queryer Queryer) ([]*Payee, error) {
	db.log.Debug("Selecting payees")

	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM payees
		WHERE valid_to_timestamp = 'infinity'
		ORDER BY id
	`, payeeColumnsStr)

	rows, err := queryer.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("selecting payees: %w", err)
	}
	defer rows.Close()
	db.log.Debugw("Selected payees", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	payees, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payee])
	if err != nil {
		return nil, fmt.Errorf("selecting payees: %w", err)
	}
	db.log.Debugw("Selected payees scanned", zap.Int("number_of_payees", len(payees)))
	return structsToPointers(payees), nil
}

func (db DB) SelectPayeesByRequestID(ctx context.Context, queryer Queryer, requestIDs ...string) (map[string]*Payee, error) {
	db.log.Debugw("Selecting payees by request ID", zap.String("request_ids", fmt.Sprintf("%+v", requestIDs)))

	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM payees
		WHERE request_id = ANY($1::TEXT[])
	`, payeeColumnsStr)

	ids := make([]pgtype.Text, 0, len(requestIDs))
	for _, id := range requestIDs {
		ids = append(ids, pgtype.Text{String: id, Valid: true})
	}

	rows, err := queryer.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("selecting payees by request ID: %w", err)
	}
	defer rows.Close()
	db.log.Debugw("Selected payees by request ID", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	payees, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payee])
	if err != nil {
		return nil, fmt.Errorf("selecting payees by request ID: %w", err)
	}
	db.log.Debugw("Selected payees by request ID scanned", zap.Int("number_of_payees", len(payees)))
	return mapByRequestID(structsToPointers(payees)), nil
}

func (db DB) SelectPayeesByID(ctx context.Context, queryer Queryer, payeeIDs ...string) (map[string]*Payee, error) {
	db.log.Debugw("Selecting payees by ID", zap.String("payee_ids", fmt.Sprintf("%+v", payeeIDs)))

	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM payees
		WHERE valid_to_timestamp = 'infinity'
		AND id = ANY($1::TEXT[])
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
	db.log.Debugw("Selected payees by ID", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	payees, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payee])
	if err != nil {
		return nil, fmt.Errorf("selecting payees by ID: %w", err)
	}
	db.log.Debugw("Selected payees by ID scanned", zap.Int("number_of_payees", len(payees)))
	return mapByID(structsToPointers(payees)), nil
}

func (db DB) SelectPayeesByName(ctx context.Context, queryer Queryer, payeeNames ...string) (map[string]*Payee, error) {
	db.log.Debugw("Selecting payees by name", zap.String("payee_names", fmt.Sprintf("%+v", payeeNames)))

	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM payees
		WHERE valid_to_timestamp = 'infinity'
		AND name = ANY($1::TEXT[])
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
	db.log.Debugw("Selected payees by ID", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	payees, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payee])
	if err != nil {
		return nil, fmt.Errorf("selecting payees by name: %w", err)
	}
	db.log.Debugw("Selected payees by name scanned", zap.Int("number_of_payees", len(payees)))
	return mapByName(structsToPointers(payees)), nil
}

func payeesToArgs(payees []*Payee) []any {
	requestIDs := make([]pgtype.Text, 0, len(payees))
	validFromTimestamps := make([]pgtype.Timestamptz, 0, len(payees))
	validToTimestamps := make([]pgtype.Timestamptz, 0, len(payees))
	ids := make([]pgtype.Text, 0, len(payees))
	names := make([]pgtype.Text, 0, len(payees))
	for _, payee := range payees {
		requestIDs = append(requestIDs, payee.RequestID)
		validFromTimestamps = append(validFromTimestamps, payee.ValidFromTimestamp)
		validToTimestamps = append(validToTimestamps, payee.ValidToTimestamp)
		ids = append(ids, payee.ID)
		names = append(names, payee.Name)
	}
	return []any{
		requestIDs,
		validFromTimestamps,
		validToTimestamps,
		ids,
		names,
	}
}
