package db_test

import (
	"context"

	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *dbSuite) TestInsertPayees() {
	ids, err := db.DB{}.InsertPayees(context.Background(), s.conn, []*db.Payee{
		{
			ID:   pgtype.Text{String: "id-1", Valid: true},
			Name: pgtype.Text{String: "name-1", Valid: true},
		},
		{
			ID:   pgtype.Text{String: "id-2", Valid: true},
			Name: pgtype.Text{String: "name-2", Valid: true},
		},
		{
			ID:   pgtype.Text{String: "id-3", Valid: true},
			Name: pgtype.Text{String: "name-3", Valid: true},
		},
	}...)
	s.NoError(err)
	s.ElementsMatch([]string{"id-1", "id-2", "id-3"}, ids)
}

func (s *dbSuite) TestSelectPayees() {
	expectedPayees := []*db.Payee{
		{
			ID:   pgtype.Text{String: "id-1", Valid: true},
			Name: pgtype.Text{String: "name-1", Valid: true},
		},
		{
			ID:   pgtype.Text{String: "id-2", Valid: true},
			Name: pgtype.Text{String: "name-2", Valid: true},
		},
		{
			ID:   pgtype.Text{String: "id-3", Valid: true},
			Name: pgtype.Text{String: "name-3", Valid: true},
		},
	}
	_, err := db.DB{}.InsertPayees(context.Background(), s.conn, expectedPayees...)
	s.Require().NoError(err)

	actualPayees, err := db.DB{}.SelectPayees(context.Background(), s.conn)
	s.NoError(err)
	s.CMPEqual(expectedPayees, actualPayees)
}

func (s *dbSuite) TestSelectPayeesByID() {
	payees := []*db.Payee{
		{
			ID:   pgtype.Text{String: "id-1", Valid: true},
			Name: pgtype.Text{String: "name-1", Valid: true},
		},
		{
			ID:   pgtype.Text{String: "id-2", Valid: true},
			Name: pgtype.Text{String: "name-2", Valid: true},
		},
		{
			ID:   pgtype.Text{String: "id-3", Valid: true},
			Name: pgtype.Text{String: "name-3", Valid: true},
		},
	}
	_, err := db.DB{}.InsertPayees(context.Background(), s.conn, payees...)
	s.Require().NoError(err)

	expectedPayees := map[string]*db.Payee{
		"id-1": payees[0],
		"id-3": payees[2],
	}
	actualPayees, err := db.DB{}.SelectPayeesByID(context.Background(), s.conn, "id-1", "id-3")
	s.NoError(err)
	s.CMPEqual(expectedPayees, actualPayees)
}

func (s *dbSuite) TestSelectPayeesByName() {
	payees := []*db.Payee{
		{
			ID:   pgtype.Text{String: "id-1", Valid: true},
			Name: pgtype.Text{String: "name-1", Valid: true},
		},
		{
			ID:   pgtype.Text{String: "id-2", Valid: true},
			Name: pgtype.Text{String: "name-2", Valid: true},
		},
		{
			ID:   pgtype.Text{String: "id-3", Valid: true},
			Name: pgtype.Text{String: "name-3", Valid: true},
		},
	}
	_, err := db.DB{}.InsertPayees(context.Background(), s.conn, payees...)
	s.Require().NoError(err)

	expectedPayees := map[string]*db.Payee{
		"name-1": payees[0],
		"name-3": payees[2],
	}
	actualPayees, err := db.DB{}.SelectPayeesByName(context.Background(), s.conn, "name-1", "name-3")
	s.NoError(err)
	s.CMPEqual(expectedPayees, actualPayees)
}
