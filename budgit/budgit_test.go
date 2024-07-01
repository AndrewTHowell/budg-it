package budgit_test

import (
	"testing"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func TestBudgit(t *testing.T) {
	suite.Run(t, new(budgitSuite))
}

type budgitSuite struct {
	suite.Suite
}

func (s *budgitSuite) TestNewBudget() {
	s.Run("ReturnsBudgetWithGeneratedUUID", func() {
		budget, err := budgit.NewBudget("", "")
		s.NoError(err, "did not expect budget to error")
		s.Require().NotEmpty(budget.ID, "expected budget to have non-empty ID")
		_, err = uuid.Parse(budget.ID)
		s.NoError(err, "expected budget to have UUID ID")
	})
	s.Run("ReturnsBudgetWithGivenName", func() {
		expectedName := "name"
		budget, err := budgit.NewBudget(expectedName, "")
		s.NoError(err, "did not expect budget to error")
		s.Equal(expectedName, budget.Name)
	})
	s.Run("ReturnsBudgetWithGivenCurrency", func() {
		expectedCurrency := "GBP"
		budget, err := budgit.NewBudget("", expectedCurrency)
		s.NoError(err, "did not expect budget to error")
		s.Equal(expectedCurrency, budget.Currency)
	})
}
