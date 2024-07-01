package budgit_test

import (
	"github.com/andrewthowell/budgit/budgit"
	"github.com/google/uuid"
)

func (s *budgitSuite) TestNewAccount() {
	s.Run("ReturnsAccountWithGeneratedUUID", func() {
		account := budgit.NewAccount("", "")
		s.Require().NotEmpty(account.ID, "expected account to have non-empty ID")
		_, err := uuid.Parse(account.ID)
		s.NoError(err, "expected account to have UUID ID")
	})
	s.Run("ReturnsAccountWithGivenBudgetID", func() {
		expectedBudgetID := uuid.UUID{1}.String()
		account := budgit.NewAccount(expectedBudgetID, "")
		s.Equal(expectedBudgetID, account.BudgetID)
	})
	s.Run("ReturnsAccountWithGivenName", func() {
		expectedName := "name"
		account := budgit.NewAccount("", expectedName)
		s.Equal(expectedName, account.Name)
	})
}
