package budgit_test

import (
	"github.com/andrewthowell/budgit/budgit"
	"github.com/google/uuid"
)

func (s *budgitSuite) TestNewAccount() {
	s.Run("ReturnsAccountWithGeneratedUUID", func() {
		account := budgit.NewAccount(uuid.UUID{}, "")
		s.NotEmpty(account.ID, "expected account to have non-empty ID")
	})
	s.Run("ReturnsAccountWithGivenBudgetID", func() {
		expectedBudgetID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		account := budgit.NewAccount(expectedBudgetID, "")
		s.Equal(expectedBudgetID, account.BudgetID)
	})
	s.Run("ReturnsAccountWithGivenName", func() {
		expectedName := "name"
		account := budgit.NewAccount(uuid.UUID{}, expectedName)
		s.Equal(expectedName, account.Name)
	})
}
