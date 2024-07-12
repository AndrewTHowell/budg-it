package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID               pgtype.Text `db:"id"`
	Name             pgtype.Text `db:"name"`
	ClearedBalance   pgtype.Int8 `db:"cleared_balance"`
	EffectiveBalance pgtype.Int8 `db:"effective_balance"`
	// Fields concerning the linked external account. Optional.
	ExternalID                pgtype.Text        `db:"external_id"`
	ExternalIntegrationID     pgtype.Text        `db:"external_integration_id"`
	ExternalLastSyncTimestamp pgtype.Timestamptz `db:"external_last_sync_timestamp"`
	ExternalClearedBalance    pgtype.Int8        `db:"external_cleared_balance"`
	ExternalEffectiveBalance  pgtype.Int8        `db:"external_effective_balance"`
}

func (a Account) GetID() string {
	return a.ID.String
}

var (
	accountColumns    = getAllDBColumns(Account{})
	accountColumnsStr = strings.Join(accountColumns, ", ")
)

func (db DB) InsertAccounts(ctx context.Context, queryer Queryer, accounts ...*Account) ([]string, error) {
	sql := fmt.Sprintf(`
		INSERT INTO accounts (%[1]s)
		(
			SELECT %[1]s
			FROM UNNEST(
				$1::TEXT[],
				$2::TEXT[],
				$3::BIGINT[],
				$4::BIGINT[],
				$5::TEXT[],
				$6::TEXT[],
				$7::TIMESTAMPTZ[],
				$8::BIGINT[],
				$9::BIGINT[]
			)
			AS u(%[1]s)
		)
		ON CONFLICT DO NOTHING
		RETURNING id;
	`, accountColumnsStr)

	rows, err := queryer.Query(ctx, sql, accountsToArgs(accounts)...)
	if err != nil {
		return nil, fmt.Errorf("inserting %d accounts: %w", len(accounts), err)
	}
	defer rows.Close()

	ids, err := rowsToIDs(rows)
	if err != nil {
		return nil, fmt.Errorf("inserting %d accounts: %w", len(accounts), err)
	}
	return ids, nil
}

func (db DB) SelectAccounts(ctx context.Context, queryer Queryer) ([]*Account, error) {
	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM accounts
		ORDER BY id
	`, accountColumnsStr)

	rows, err := queryer.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("selecting accounts: %w", err)
	}
	defer rows.Close()

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Account])
	if err != nil {
		return nil, fmt.Errorf("selecting accounts: %w", err)
	}
	return structsToPointers(accounts), nil
}

func (db DB) SelectAccountsByID(ctx context.Context, queryer Queryer, accountIDs ...string) (map[string]*Account, error) {
	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM accounts
		WHERE id = ANY($1::TEXT[])
	`, accountColumnsStr)

	ids := make([]pgtype.Text, 0, len(accountIDs))
	for _, id := range accountIDs {
		ids = append(ids, pgtype.Text{String: id, Valid: true})
	}

	rows, err := queryer.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("selecting accounts by ID: %w", err)
	}
	defer rows.Close()

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Account])
	if err != nil {
		return nil, fmt.Errorf("selecting accounts by ID: %w", err)
	}
	return mapByID(structsToPointers(accounts)), nil
}

func accountsToArgs(accounts []*Account) []any {
	ids := make([]pgtype.Text, 0, len(accounts))
	names := make([]pgtype.Text, 0, len(accounts))
	cleared_balances := make([]pgtype.Int8, 0, len(accounts))
	effective_balances := make([]pgtype.Int8, 0, len(accounts))
	external_ids := make([]pgtype.Text, 0, len(accounts))
	external_integration_ids := make([]pgtype.Text, 0, len(accounts))
	external_last_sync_timestamp := make([]pgtype.Timestamptz, 0, len(accounts))
	external_cleared_balance := make([]pgtype.Int8, 0, len(accounts))
	external_effective_balance := make([]pgtype.Int8, 0, len(accounts))
	for _, account := range accounts {
		ids = append(ids, account.ID)
		names = append(names, account.Name)
		cleared_balances = append(cleared_balances, account.ClearedBalance)
		effective_balances = append(effective_balances, account.EffectiveBalance)
		external_ids = append(external_ids, account.ExternalID)
		external_integration_ids = append(external_integration_ids, account.ExternalIntegrationID)
		external_last_sync_timestamp = append(external_last_sync_timestamp, account.ExternalLastSyncTimestamp)
		external_cleared_balance = append(external_cleared_balance, account.ExternalClearedBalance)
		external_effective_balance = append(external_effective_balance, account.ExternalEffectiveBalance)
	}
	return []any{
		ids,
		names,
		cleared_balances,
		effective_balances,
		external_ids,
		external_integration_ids,
		external_last_sync_timestamp,
		external_cleared_balance,
		external_effective_balance,
	}
}
