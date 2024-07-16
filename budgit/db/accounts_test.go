package db_test

import (
	"context"
	"time"

	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *dbSuite) TestInsertAccounts() {
	ids, err := s.db.InsertAccounts(context.Background(), s.conn, []*db.Account{
		{
			RequestID:                 pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-1", Valid: true},
			Name:                      pgtype.Text{String: "name-1", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 1, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 2, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-1", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-1", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-1", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(3, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 3, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 4, Valid: true},
		},
		{
			RequestID:                 pgtype.Text{String: "request_id-2", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(4, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-2", Valid: true},
			Name:                      pgtype.Text{String: "name-2", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 2, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 3, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-2", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-2", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-2", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(6, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 4, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 5, Valid: true},
		},
		{
			RequestID:                 pgtype.Text{String: "request_id-3", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(7, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-3", Valid: true},
			Name:                      pgtype.Text{String: "name-3", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 3, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 4, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-3", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-3", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-3", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(9, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 5, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 6, Valid: true},
		},
	}...)
	s.NoError(err)
	s.ElementsMatch([]string{"id-1", "id-2", "id-3"}, ids)
}

func (s *dbSuite) TestUpdateAccountValidToTimestamps() {
	_, err := s.db.InsertAccounts(context.Background(), s.conn, []*db.Account{
		{
			RequestID:                 pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-1", Valid: true},
			Name:                      pgtype.Text{String: "name-1", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 1, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 2, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-1", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-1", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-1", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 3, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 4, Valid: true},
		},
		{
			RequestID:                 pgtype.Text{String: "request_id-2", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(3, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-2", Valid: true},
			Name:                      pgtype.Text{String: "name-2", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 2, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 3, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-2", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-2", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-2", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(4, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 4, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 5, Valid: true},
		},
		{
			RequestID:                 pgtype.Text{String: "request_id-3", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(5, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-3", Valid: true},
			Name:                      pgtype.Text{String: "name-3", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 3, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 4, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-3", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-3", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-3", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(6, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 5, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 6, Valid: true},
		},
	}...)
	s.Require().NoError(err)

	ids, err := s.db.UpdateAccountValidToTimestamps(context.Background(), s.conn, []db.ValidToTimestampUpdate{
		{
			ID:               pgtype.Text{String: "id-1", Valid: true},
			ValidToTimestamp: pgtype.Timestamptz{Time: time.Unix(10, 0).UTC(), Valid: true},
		},
		{
			ID:               pgtype.Text{String: "id-3", Valid: true},
			ValidToTimestamp: pgtype.Timestamptz{Time: time.Unix(13, 0).UTC(), Valid: true},
		},
	}...)
	s.NoError(err)
	s.ElementsMatch([]string{"id-1", "id-3"}, ids)

	s.Run("AccountsUpdatedInDB", func() {
		actualAccounts, err := s.db.SelectAccountsByRequestID(context.Background(), s.conn, "request_id-1", "request_id-2", "request_id-3")
		s.NoError(err)
		s.CMPEqual(map[string]*db.Account{
			"request_id-1": {
				RequestID:                 pgtype.Text{String: "request_id-1", Valid: true},
				ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
				ValidToTimestamp:          pgtype.Timestamptz{Time: time.Unix(10, 0).UTC(), Valid: true},
				ID:                        pgtype.Text{String: "id-1", Valid: true},
				Name:                      pgtype.Text{String: "name-1", Valid: true},
				ClearedBalance:            pgtype.Int8{Int64: 1, Valid: true},
				EffectiveBalance:          pgtype.Int8{Int64: 2, Valid: true},
				ExternalID:                pgtype.Text{String: "external_id-1", Valid: true},
				ExternalName:              pgtype.Text{String: "external_name-1", Valid: true},
				ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-1", Valid: true},
				ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
				ExternalClearedBalance:    pgtype.Int8{Int64: 3, Valid: true},
				ExternalEffectiveBalance:  pgtype.Int8{Int64: 4, Valid: true},
			},
			"request_id-2": {
				RequestID:                 pgtype.Text{String: "request_id-2", Valid: true},
				ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(3, 0).UTC(), Valid: true},
				ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
				ID:                        pgtype.Text{String: "id-2", Valid: true},
				Name:                      pgtype.Text{String: "name-2", Valid: true},
				ClearedBalance:            pgtype.Int8{Int64: 2, Valid: true},
				EffectiveBalance:          pgtype.Int8{Int64: 3, Valid: true},
				ExternalID:                pgtype.Text{String: "external_id-2", Valid: true},
				ExternalName:              pgtype.Text{String: "external_name-2", Valid: true},
				ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-2", Valid: true},
				ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(4, 0).UTC(), Valid: true},
				ExternalClearedBalance:    pgtype.Int8{Int64: 4, Valid: true},
				ExternalEffectiveBalance:  pgtype.Int8{Int64: 5, Valid: true},
			},
			"request_id-3": {
				RequestID:                 pgtype.Text{String: "request_id-3", Valid: true},
				ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(5, 0).UTC(), Valid: true},
				ValidToTimestamp:          pgtype.Timestamptz{Time: time.Unix(13, 0).UTC(), Valid: true},
				ID:                        pgtype.Text{String: "id-3", Valid: true},
				Name:                      pgtype.Text{String: "name-3", Valid: true},
				ClearedBalance:            pgtype.Int8{Int64: 3, Valid: true},
				EffectiveBalance:          pgtype.Int8{Int64: 4, Valid: true},
				ExternalID:                pgtype.Text{String: "external_id-3", Valid: true},
				ExternalName:              pgtype.Text{String: "external_name-3", Valid: true},
				ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-3", Valid: true},
				ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(6, 0).UTC(), Valid: true},
				ExternalClearedBalance:    pgtype.Int8{Int64: 5, Valid: true},
				ExternalEffectiveBalance:  pgtype.Int8{Int64: 6, Valid: true},
			},
		}, actualAccounts)
	})
}

func (s *dbSuite) TestSelectAccounts() {
	expectedAccounts := []*db.Account{
		{
			RequestID:                 pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-1", Valid: true},
			Name:                      pgtype.Text{String: "name-1", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 1, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 2, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-1", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-1", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-1", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 3, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 4, Valid: true},
		},
	}
	_, err := s.db.InsertAccounts(context.Background(), s.conn, expectedAccounts...)
	s.Require().NoError(err)

	actualAccounts, err := s.db.SelectAccounts(context.Background(), s.conn)
	s.NoError(err)
	s.CMPEqual(expectedAccounts, actualAccounts)
}

func (s *dbSuite) TestSelectAccountsByRequestID() {
	accounts := []*db.Account{
		{
			RequestID:                 pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-1", Valid: true},
			Name:                      pgtype.Text{String: "name-1", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 1, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 2, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-1", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-1", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-1", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 3, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 4, Valid: true},
		},
		{
			RequestID:                 pgtype.Text{String: "request_id-2", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(3, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-2", Valid: true},
			Name:                      pgtype.Text{String: "name-2", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 2, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 3, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-2", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-2", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-2", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(4, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 4, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 5, Valid: true},
		},
		{
			RequestID:                 pgtype.Text{String: "request_id-3", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(5, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-3", Valid: true},
			Name:                      pgtype.Text{String: "name-3", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 3, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 4, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-3", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-3", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-3", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(6, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 5, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 6, Valid: true},
		},
	}
	_, err := s.db.InsertAccounts(context.Background(), s.conn, accounts...)
	s.Require().NoError(err)

	expectedAccounts := map[string]*db.Account{
		"request_id-1": accounts[0],
		"request_id-3": accounts[2],
	}
	actualAccounts, err := s.db.SelectAccountsByRequestID(context.Background(), s.conn, "request_id-1", "request_id-3")
	s.NoError(err)
	s.CMPEqual(expectedAccounts, actualAccounts)
}

func (s *dbSuite) TestSelectAccountsByID() {
	accounts := []*db.Account{
		{
			RequestID:                 pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-1", Valid: true},
			Name:                      pgtype.Text{String: "name-1", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 1, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 2, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-1", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-1", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-1", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 3, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 4, Valid: true},
		},
		{
			RequestID:                 pgtype.Text{String: "request_id-2", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(3, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-2", Valid: true},
			Name:                      pgtype.Text{String: "name-2", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 2, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 3, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-2", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-2", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-2", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(4, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 4, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 5, Valid: true},
		},
		{
			RequestID:                 pgtype.Text{String: "request_id-3", Valid: true},
			ValidFromTimestamp:        pgtype.Timestamptz{Time: time.Unix(5, 0).UTC(), Valid: true},
			ValidToTimestamp:          pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                        pgtype.Text{String: "id-3", Valid: true},
			Name:                      pgtype.Text{String: "name-3", Valid: true},
			ClearedBalance:            pgtype.Int8{Int64: 3, Valid: true},
			EffectiveBalance:          pgtype.Int8{Int64: 4, Valid: true},
			ExternalID:                pgtype.Text{String: "external_id-3", Valid: true},
			ExternalName:              pgtype.Text{String: "external_name-3", Valid: true},
			ExternalIntegrationID:     pgtype.Text{String: "external_integration_id-3", Valid: true},
			ExternalLastSyncTimestamp: pgtype.Timestamptz{Time: time.Unix(6, 0).UTC(), Valid: true},
			ExternalClearedBalance:    pgtype.Int8{Int64: 5, Valid: true},
			ExternalEffectiveBalance:  pgtype.Int8{Int64: 6, Valid: true},
		},
	}
	_, err := s.db.InsertAccounts(context.Background(), s.conn, accounts...)
	s.Require().NoError(err)

	expectedAccounts := map[string]*db.Account{
		"id-1": accounts[0],
		"id-3": accounts[2],
	}
	actualAccounts, err := s.db.SelectAccountsByID(context.Background(), s.conn, "id-1", "id-3")
	s.NoError(err)
	s.CMPEqual(expectedAccounts, actualAccounts)
}
