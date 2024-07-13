package budgit_test

import (
	"time"

	"github.com/andrewthowell/budgit/budgit"
)

func (s *budgitSuite) TestTransaction() {
	testCases := []struct {
		name                           string
		mirrorID                       string
		transaction, mirrorTransaction *budgit.Transaction
	}{
		{
			name:        "EmptyTransaction",
			mirrorID:    "",
			transaction: &budgit.Transaction{},
			mirrorTransaction: &budgit.Transaction{
				IsPayeeInternal: true,
			},
		},
		{
			name:     "PopulatedTransaction",
			mirrorID: "mirror_id-1",
			transaction: &budgit.Transaction{
				ID:              "id-1",
				EffectiveDate:   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
				AccountID:       "account_id-1",
				PayeeID:         "payee_id-1",
				IsPayeeInternal: true,
				Amount:          1,
				Cleared:         true,
			},
			mirrorTransaction: &budgit.Transaction{
				ID:              "mirror_id-1",
				EffectiveDate:   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
				AccountID:       "payee_id-1",
				PayeeID:         "account_id-1",
				IsPayeeInternal: false,
				Amount:          -1,
				Cleared:         true,
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Run("MirroredTransactionEqualsMirrorTransaction", func() {
				s.CMPEqual(tc.mirrorTransaction, tc.transaction.Mirror(tc.mirrorID))
			})
			s.Run("MirroredMirrorTransactionEqualsTransaction", func() {
				s.CMPEqual(tc.transaction, tc.mirrorTransaction.Mirror(tc.transaction.ID))
			})
			s.Run("TransactionToMirrorToTransaction", func() {
				s.CMPEqual(tc.transaction, tc.transaction.Mirror(tc.mirrorID).Mirror(tc.transaction.ID))
			})
			s.Run("MirrorToTransactionToMirror", func() {
				s.CMPEqual(tc.mirrorTransaction, tc.mirrorTransaction.Mirror(tc.transaction.ID).Mirror(tc.mirrorID))
			})
		})
	}
}
