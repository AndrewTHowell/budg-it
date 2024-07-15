package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type Transaction struct {
	RequestID          pgtype.Text        `db:"request_id"`
	ValidFromTimestamp pgtype.Timestamptz `db:"valid_from_timestamp"`
	ValidToTimestamp   pgtype.Timestamptz `db:"valid_to_timestamp"`
	ID                 pgtype.Text        `db:"id"`
	EffectiveDate      pgtype.Date        `db:"effective_date"`
	AccountID          pgtype.Text        `db:"account_id"`
	PayeeID            pgtype.Text        `db:"payee_id"`
	IsPayeeInternal    pgtype.Bool        `db:"is_payee_internal"`
	Amount             pgtype.Int8        `db:"amount"`
	Cleared            pgtype.Bool        `db:"cleared"`
}

func (p Transaction) GetID() string {
	return p.ID.String
}

func (p Transaction) GetRequestID() string {
	return p.RequestID.String
}

func (p Transaction) GetAccountID() string {
	return p.AccountID.String
}

var (
	transactionColumns    = getAllDBColumns(Transaction{})
	transactionColumnsStr = strings.Join(transactionColumns, ", ")
)

func (db DB) InsertTransactions(ctx context.Context, queryer Queryer, transactions ...*Transaction) ([]string, error) {
	db.log.Debugw("Inserting transactions", zap.Int("number_of_transactions", len(transactions)))

	sql := fmt.Sprintf(`
		INSERT INTO transactions (%[1]s)
		(
			SELECT %[1]s
			FROM UNNEST(
				$1::TEXT[],
				$2::TIMESTAMPTZ[],
				$3::TIMESTAMPTZ[],
				$4::TEXT[],
				$5::DATE[],
				$6::TEXT[],
				$7::TEXT[],
				$8::BOOL[],
				$9::BIGINT[],
				$10::BOOL[]
			)
			AS u(%[1]s)
		)
		ON CONFLICT DO NOTHING
		RETURNING id;
	`, transactionColumnsStr)

	rows, err := queryer.Query(ctx, sql, transactionsToArgs(transactions)...)
	if err != nil {
		return nil, fmt.Errorf("inserting %d transactions: %w", len(transactions), err)
	}
	defer rows.Close()
	db.log.Debugw("Inserted transactions", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	ids, err := rowsToIDs(rows)
	if err != nil {
		return nil, fmt.Errorf("inserting %d transactions: %w", len(transactions), err)
	}
	db.log.Debugw("Inserted transactions scanned", zap.String("inserted_ids", fmt.Sprintf("%v", ids)))
	return ids, nil
}

func (db DB) SelectTransactions(ctx context.Context, queryer Queryer) ([]*Transaction, error) {
	db.log.Debug("Selecting transactions")

	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM transactions
		WHERE valid_to_timestamp = 'infinity'
		ORDER BY effective_date, amount
	`, transactionColumnsStr)

	rows, err := queryer.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("selecting transactions: %w", err)
	}
	defer rows.Close()
	db.log.Debugw("Selected transactions", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	transactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		return nil, fmt.Errorf("selecting transactions: %w", err)
	}
	db.log.Debugw("Selected transactions scanned", zap.Int("number_of_transactions", len(transactions)))
	return structsToPointers(transactions), nil
}

func (db DB) SelectTransactionsByAccount(ctx context.Context, queryer Queryer, accountID string) ([]*Transaction, error) {
	db.log.Debug("Selecting transactions by account")

	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM transactions
		WHERE valid_to_timestamp = 'infinity'
		AND account_id = $1
		ORDER BY effective_date, amount
	`, transactionColumnsStr)

	rows, err := queryer.Query(ctx, sql, pgtype.Text{String: accountID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("selecting transactions by name: %w", err)
	}
	defer rows.Close()
	db.log.Debugw("Selected transactions", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	transactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		return nil, fmt.Errorf("selecting transactions by name: %w", err)
	}
	db.log.Debugw("Selected transactions scanned", zap.Int("number_of_transactions", len(transactions)))
	return structsToPointers(transactions), nil
}

func (db DB) SelectTransactionsByRequestID(ctx context.Context, queryer Queryer, requestIDs ...string) (map[string]*Transaction, error) {
	db.log.Debugw("Selecting transactions by request ID", zap.String("request_ids", fmt.Sprintf("%+v", requestIDs)))

	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM transactions
		WHERE valid_to_timestamp = 'infinity'
		AND request_id = ANY($1::TEXT[])
	`, transactionColumnsStr)

	ids := make([]pgtype.Text, 0, len(requestIDs))
	for _, id := range requestIDs {
		ids = append(ids, pgtype.Text{String: id, Valid: true})
	}

	rows, err := queryer.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("selecting transactions by request ID: %w", err)
	}
	defer rows.Close()
	db.log.Debugw("Selected transactions by request ID", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	transactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		return nil, fmt.Errorf("selecting transactions by request ID: %w", err)
	}
	db.log.Debugw("Selected transactions by request ID scanned", zap.Int("number_of_transactions", len(transactions)))
	return mapByRequestID(structsToPointers(transactions)), nil
}

func (db DB) SelectTransactionsByID(ctx context.Context, queryer Queryer, transactionIDs ...string) (map[string]*Transaction, error) {
	db.log.Debugw("Selecting transactions by ID", zap.String("transaction_ids", fmt.Sprintf("%+v", transactionIDs)))

	sql := fmt.Sprintf(`
		SELECT %[1]s
		FROM transactions
		WHERE valid_to_timestamp = 'infinity'
		AND id = ANY($1::TEXT[])
	`, transactionColumnsStr)

	ids := make([]pgtype.Text, 0, len(transactionIDs))
	for _, id := range transactionIDs {
		ids = append(ids, pgtype.Text{String: id, Valid: true})
	}

	rows, err := queryer.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("selecting transactions by ID: %w", err)
	}
	defer rows.Close()
	db.log.Debugw("Selected transactions by ID", zap.Int64("rows_affected", rows.CommandTag().RowsAffected()))

	transactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		return nil, fmt.Errorf("selecting transactions by ID: %w", err)
	}
	db.log.Debugw("Selected transactions by ID scanned", zap.Int("number_of_transactions", len(transactions)))
	return mapByID(structsToPointers(transactions)), nil
}

func transactionsToArgs(transactions []*Transaction) []any {
	requestIDs := make([]pgtype.Text, 0, len(transactions))
	validFromTimestamps := make([]pgtype.Timestamptz, 0, len(transactions))
	validToTimestamps := make([]pgtype.Timestamptz, 0, len(transactions))
	ids := make([]pgtype.Text, 0, len(transactions))
	effective_dates := make([]pgtype.Date, 0, len(transactions))
	account_ids := make([]pgtype.Text, 0, len(transactions))
	payee_ids := make([]pgtype.Text, 0, len(transactions))
	is_payee_internals := make([]pgtype.Bool, 0, len(transactions))
	amounts := make([]pgtype.Int8, 0, len(transactions))
	cleareds := make([]pgtype.Bool, 0, len(transactions))
	for _, transaction := range transactions {
		requestIDs = append(requestIDs, transaction.RequestID)
		validFromTimestamps = append(validFromTimestamps, transaction.ValidFromTimestamp)
		validToTimestamps = append(validToTimestamps, transaction.ValidToTimestamp)
		ids = append(ids, transaction.ID)
		effective_dates = append(effective_dates, transaction.EffectiveDate)
		account_ids = append(account_ids, transaction.AccountID)
		payee_ids = append(payee_ids, transaction.PayeeID)
		is_payee_internals = append(is_payee_internals, transaction.IsPayeeInternal)
		amounts = append(amounts, transaction.Amount)
		cleareds = append(cleareds, transaction.Cleared)
	}
	return []any{
		requestIDs,
		validFromTimestamps,
		validToTimestamps,
		ids,
		effective_dates,
		account_ids,
		payee_ids,
		is_payee_internals,
		amounts,
		cleareds,
	}
}
