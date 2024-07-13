package dbconvert_test

import (
	"time"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/db/dbconvert"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *convertSuite) TestTransaction() {
	testCases := []struct {
		name              string
		dbTransaction     *db.Transaction
		budgitTransaction *budgit.Transaction
	}{
		{
			name:              "EmptyTransaction",
			dbTransaction:     &db.Transaction{},
			budgitTransaction: &budgit.Transaction{},
		},
		{
			name: "PopulatedTransaction",
			dbTransaction: &db.Transaction{
				ID:              pgtype.Text{String: "id-1", Valid: true},
				EffectiveDate:   pgtype.Date{Time: time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC), Valid: true},
				AccountID:       pgtype.Text{String: "account_id-1", Valid: true},
				PayeeID:         pgtype.Text{String: "payee_id-1", Valid: true},
				IsPayeeInternal: pgtype.Bool{Bool: true, Valid: true},
				Amount:          pgtype.Int8{Int64: 1, Valid: true},
				Cleared:         pgtype.Bool{Bool: true, Valid: true},
			},
			budgitTransaction: &budgit.Transaction{
				ID:              "id-1",
				EffectiveDate:   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
				AccountID:       "account_id-1",
				PayeeID:         "payee_id-1",
				IsPayeeInternal: true,
				Amount:          1,
				Cleared:         true,
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Run("ToTransaction", func() {
				s.CMPEqual(tc.budgitTransaction, dbconvert.ToTransactions(tc.dbTransaction))
			})
			s.Run("FromTransaction", func() {
				s.CMPEqual(tc.dbTransaction, dbconvert.FromTransactions(tc.budgitTransaction))
			})
			s.Run("FromTransactionToTransaction", func() {
				s.CMPEqual(tc.dbTransaction, dbconvert.FromTransactions(dbconvert.ToTransactions(tc.dbTransaction)...))
			})
			s.Run("ToTransactionFromTransaction", func() {
				s.CMPEqual(tc.budgitTransaction, dbconvert.ToTransactions(dbconvert.FromTransactions(tc.budgitTransaction)...))
			})
		})
	}
}
