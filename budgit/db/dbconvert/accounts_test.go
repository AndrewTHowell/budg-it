package dbconvert_test

import (
	"time"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/db/dbconvert"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *convertSuite) TestAccount() {
	testCases := []struct {
		name          string
		dbAccount     *db.Account
		budgitAccount *budgit.Account
	}{
		{
			name:          "EmptyAccount",
			dbAccount:     &db.Account{},
			budgitAccount: &budgit.Account{},
		},
		{
			name: "PopulatedAccount",
			dbAccount: &db.Account{
				ID:                        pgtype.Text{String: "id-1", Valid: true},
				Name:                      pgtype.Text{String: "name-1", Valid: true},
				ClearedBalance:            pgtype.Int8{Int64: 1, Valid: true},
				EffectiveBalance:          pgtype.Int8{Int64: 2, Valid: true},
				ExternalID:                pgtype.Text{String: "external_id-1", Valid: true},
				ExternalName:              pgtype.Text{String: "external_name-1", Valid: true},
				ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-1", Valid: true},
				ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
				ExternalClearedBalance:    pgtype.Int8{Int64: 3, Valid: true},
				ExternalEffectiveBalance:  pgtype.Int8{Int64: 4, Valid: true},
			},
			budgitAccount: &budgit.Account{
				ID:   "id-1",
				Name: "name-1",
				Balance: budgit.Balance{
					ClearedBalance:   1,
					EffectiveBalance: 2,
				},
				ExternalAccount: &budgit.ExternalAccount{
					ID:                "external_id-1",
					Name:              "external_name-1",
					IntegrationID:     "external_integration_id-1",
					LastSyncTimestamp: time.Unix(1, 0).UTC(),
					Balance: budgit.Balance{
						ClearedBalance:   3,
						EffectiveBalance: 4,
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Run("ToAccount", func() {
				s.CMPEqual(tc.budgitAccount, dbconvert.ToAccounts(tc.dbAccount)[0])
			})
			s.Run("FromAccount", func() {
				s.CMPEqual(tc.dbAccount, dbconvert.FromAccounts(tc.budgitAccount)[0])
			})
			s.Run("FromAccountToAccount", func() {
				s.CMPEqual(tc.dbAccount, dbconvert.FromAccounts(dbconvert.ToAccounts(tc.dbAccount)...)[0])
			})
			s.Run("ToAccountFromAccount", func() {
				s.CMPEqual(tc.budgitAccount, dbconvert.ToAccounts(dbconvert.FromAccounts(tc.budgitAccount)...)[0])
			})
		})
	}
}
