package db_test

import (
	"context"
	"time"

	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *dbSuite) TestInsertTransactions() {
	ids, err := db.DB{}.InsertTransactions(context.Background(), s.conn, []*db.Transaction{
		{
			ID:              pgtype.Text{String: "id-1", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-1", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-1", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 1, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:              pgtype.Text{String: "id-2", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-2", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-2", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 2, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:              pgtype.Text{String: "id-3", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-3", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-3", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 3, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
	}...)
	s.NoError(err)
	s.ElementsMatch([]string{"id-1", "id-2", "id-3"}, ids)
}

func (s *dbSuite) TestSelectTransactions() {
	expectedTransactions := []*db.Transaction{
		{
			ID:              pgtype.Text{String: "id-1", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-1", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-1", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 1, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:              pgtype.Text{String: "id-2", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-2", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-2", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 2, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:              pgtype.Text{String: "id-3", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-3", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-3", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 3, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
	}
	_, err := db.DB{}.InsertTransactions(context.Background(), s.conn, expectedTransactions...)
	s.Require().NoError(err)

	actualTransactions, err := db.DB{}.SelectTransactions(context.Background(), s.conn)
	s.NoError(err)
	s.CMPEqual(expectedTransactions, actualTransactions)
}

func (s *dbSuite) TestSelectTransactionsByID() {
	transactions := []*db.Transaction{
		{
			ID:              pgtype.Text{String: "id-1", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-1", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-1", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 1, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:              pgtype.Text{String: "id-2", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-2", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-2", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 2, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:              pgtype.Text{String: "id-3", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-3", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-3", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 3, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
	}
	_, err := db.DB{}.InsertTransactions(context.Background(), s.conn, transactions...)
	s.Require().NoError(err)

	expectedTransactions := map[string]*db.Transaction{
		"id-1": transactions[0],
		"id-3": transactions[2],
	}
	actualTransactions, err := db.DB{}.SelectTransactionsByID(context.Background(), s.conn, "id-1", "id-3")
	s.NoError(err)
	s.CMPEqual(expectedTransactions, actualTransactions)
}

func (s *dbSuite) TestSelectTransactionsByAccount() {
	transactions := []*db.Transaction{
		{
			ID:              pgtype.Text{String: "id-1", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-1", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-1", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 1, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:              pgtype.Text{String: "id-2", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-2", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-2", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 2, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:              pgtype.Text{String: "id-3", Valid: true},
			EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC), Valid: true},
			AccountID:       pgtype.Text{String: "account-id-1", Valid: true},
			PayeeID:         pgtype.Text{String: "payee-id-3", Valid: true},
			IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
			Amount:          pgtype.Int8{Int64: 3, Valid: true},
			Cleared:         pgtype.Bool{Bool: true, Valid: true},
		},
	}
	_, err := db.DB{}.InsertTransactions(context.Background(), s.conn, transactions...)
	s.Require().NoError(err)

	expectedTransactions := []*db.Transaction{
		transactions[0],
		transactions[2],
	}
	actualTransactions, err := db.DB{}.SelectTransactionsByAccount(context.Background(), s.conn, "account-id-1")
	s.NoError(err)
	s.CMPEqual(expectedTransactions, actualTransactions)
}
