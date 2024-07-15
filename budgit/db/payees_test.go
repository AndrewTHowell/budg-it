package db_test

import (
	"context"
	"time"

	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *dbSuite) TestInsertPayees() {
	ids, err := s.db.InsertPayees(context.Background(), s.conn, []*db.Payee{
		{
			RequestID:          pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-1", Valid: true},
			Name:               pgtype.Text{String: "name-1", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-2", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-2", Valid: true},
			Name:               pgtype.Text{String: "name-2", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-3", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-3", Valid: true},
			Name:               pgtype.Text{String: "name-3", Valid: true},
		},
	}...)
	s.NoError(err)
	s.ElementsMatch([]string{"id-1", "id-2", "id-3"}, ids)
}

func (s *dbSuite) TestSelectPayees() {
	expectedPayees := []*db.Payee{
		{
			RequestID:          pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-1", Valid: true},
			Name:               pgtype.Text{String: "name-1", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-2", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-2", Valid: true},
			Name:               pgtype.Text{String: "name-2", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-3", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(3, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-3", Valid: true},
			Name:               pgtype.Text{String: "name-3", Valid: true},
		},
	}
	_, err := s.db.InsertPayees(context.Background(), s.conn, expectedPayees...)
	s.Require().NoError(err)

	actualPayees, err := s.db.SelectPayees(context.Background(), s.conn)
	s.NoError(err)
	s.CMPEqual(expectedPayees, actualPayees)
}

func (s *dbSuite) TestSelectPayeesByRequestID() {
	payees := []*db.Payee{
		{
			RequestID:          pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-1", Valid: true},
			Name:               pgtype.Text{String: "name-1", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-2", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-2", Valid: true},
			Name:               pgtype.Text{String: "name-2", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-3", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(3, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-3", Valid: true},
			Name:               pgtype.Text{String: "name-3", Valid: true},
		},
	}
	_, err := s.db.InsertPayees(context.Background(), s.conn, payees...)
	s.Require().NoError(err)

	expectedPayees := map[string]*db.Payee{
		"request_id-1": payees[0],
		"request_id-3": payees[2],
	}
	actualPayees, err := s.db.SelectPayeesByRequestID(context.Background(), s.conn, "request_id-1", "request_id-3")
	s.NoError(err)
	s.CMPEqual(expectedPayees, actualPayees)
}

func (s *dbSuite) TestSelectPayeesByID() {
	payees := []*db.Payee{
		{
			RequestID:          pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-1", Valid: true},
			Name:               pgtype.Text{String: "name-1", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-2", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-2", Valid: true},
			Name:               pgtype.Text{String: "name-2", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-3", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(3, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-3", Valid: true},
			Name:               pgtype.Text{String: "name-3", Valid: true},
		},
	}
	_, err := s.db.InsertPayees(context.Background(), s.conn, payees...)
	s.Require().NoError(err)

	expectedPayees := map[string]*db.Payee{
		"id-1": payees[0],
		"id-3": payees[2],
	}
	actualPayees, err := s.db.SelectPayeesByID(context.Background(), s.conn, "id-1", "id-3")
	s.NoError(err)
	s.CMPEqual(expectedPayees, actualPayees)
}

func (s *dbSuite) TestSelectPayeesByName() {
	payees := []*db.Payee{
		{
			RequestID:          pgtype.Text{String: "request_id-1", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(1, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-1", Valid: true},
			Name:               pgtype.Text{String: "name-1", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-2", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(2, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-2", Valid: true},
			Name:               pgtype.Text{String: "name-2", Valid: true},
		},
		{
			RequestID:          pgtype.Text{String: "request_id-3", Valid: true},
			ValidFromTimestamp: pgtype.Timestamptz{Time: time.Unix(3, 0).UTC(), Valid: true},
			ValidToTimestamp:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
			ID:                 pgtype.Text{String: "id-3", Valid: true},
			Name:               pgtype.Text{String: "name-3", Valid: true},
		},
	}
	_, err := s.db.InsertPayees(context.Background(), s.conn, payees...)
	s.Require().NoError(err)

	expectedPayees := map[string]*db.Payee{
		"name-1": payees[0],
		"name-3": payees[2],
	}
	actualPayees, err := s.db.SelectPayeesByName(context.Background(), s.conn, "name-1", "name-3")
	s.NoError(err)
	s.CMPEqual(expectedPayees, actualPayees)
}
