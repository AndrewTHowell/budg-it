package budgit_test

import (
	date "cloud.google.com/go/civil"
	"github.com/andrewthowell/budgit/budgit"
	"github.com/google/uuid"
)

func (s *budgitSuite) TestNewTransaction() {
	s.Run("ReturnsTransactionWithGeneratedUUID", func() {
		transaction := budgit.NewTransaction("", "", date.Date{}, "", "", 0)
		s.Require().NotEmpty(transaction.ID, "expected transaction to have non-empty ID")
		_, err := uuid.Parse(transaction.ID)
		s.NoError(err, "expected transaction to have UUID ID")
	})
	s.Run("ReturnsTransactionWithGivenBudgetID", func() {
		expectedBudgetID := "budgetID"
		transaction := budgit.NewTransaction(expectedBudgetID, "", date.Date{}, "", "", 0)
		s.Equal(expectedBudgetID, transaction.BudgetID)
	})
	s.Run("ReturnsTransactionWithGivenAccountID", func() {
		expectedAccountID := "accountID"
		transaction := budgit.NewTransaction("", expectedAccountID, date.Date{}, "", "", 0)
		s.Equal(expectedAccountID, transaction.AccountID)
	})
	s.Run("ReturnsTransactionWithGivenDate", func() {
		expectedDate, _ := date.ParseDate("1970/01/01")
		transaction := budgit.NewTransaction("", "", expectedDate, "", "", 0)
		s.Equal(expectedDate, transaction.EffectiveDate)
	})
	s.Run("ReturnsTransactionWithGivenPayeeID", func() {
		expectedPayeeID := "payeeID"
		transaction := budgit.NewTransaction("", "", date.Date{}, expectedPayeeID, "", 0)
		s.Equal(expectedPayeeID, transaction.PayeeID)
	})
	s.Run("ReturnsTransactionWithGivenCategoryID", func() {
		expectedCategoryID := "categoryID"
		transaction := budgit.NewTransaction("", "", date.Date{}, "", expectedCategoryID, 0)
		s.Equal(expectedCategoryID, transaction.CategoryID)
	})
	s.Run("ReturnsTransactionWithGivenPositiveAmount", func() {
		expectedAmount := 10
		transaction := budgit.NewTransaction("", "", date.Date{}, "", "", expectedAmount)
		s.Equal(expectedAmount, transaction.Amount)
	})
	s.Run("ReturnsTransactionWithGivenNegativeAmount", func() {
		expectedAmount := -10
		transaction := budgit.NewTransaction("", "", date.Date{}, "", "", expectedAmount)
		s.Equal(expectedAmount, transaction.Amount)
	})
}
