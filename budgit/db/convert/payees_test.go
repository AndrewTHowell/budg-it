package convert_test

import (
	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/db/convert"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *convertSuite) TestPayee() {
	testCases := []struct {
		name        string
		dbPayee     *db.Payee
		budgitPayee *budgit.Payee
	}{
		{
			name:        "EmptyPayee",
			dbPayee:     &db.Payee{},
			budgitPayee: &budgit.Payee{},
		},
		{
			name: "PopulatedPayee",
			dbPayee: &db.Payee{
				ID:   pgtype.Text{String: "id-1", Valid: true},
				Name: pgtype.Text{String: "name-1", Valid: true},
			},
			budgitPayee: &budgit.Payee{
				ID:   "id-1",
				Name: "name-1",
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Run("ToPayee", func() {
				s.CMPEqual(tc.budgitPayee, convert.ToPayee(tc.dbPayee))
			})
			s.Run("FromPayee", func() {
				s.CMPEqual(tc.dbPayee, convert.FromPayee(tc.budgitPayee))
			})
			s.Run("FromPayeeToPayee", func() {
				s.CMPEqual(tc.dbPayee, convert.FromPayee(convert.ToPayee(tc.dbPayee)))
			})
			s.Run("ToPayeeFromPayee", func() {
				s.CMPEqual(tc.budgitPayee, convert.ToPayee(convert.FromPayee(tc.budgitPayee)))
			})
		})
	}
}
