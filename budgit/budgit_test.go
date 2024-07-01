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
		budget := budgit.NewBudget("", "")
		s.Require().NotEmpty(budget.ID, "expected budget to have non-empty ID")
		_, err := uuid.Parse(budget.ID.String())
		s.NoError(err, "expected budget to have UUID ID")
	})
	s.Run("ReturnsBudgetWithGivenName", func() {
		expectedName := "name"
		budget := budgit.NewBudget(expectedName, "")
		s.Equal(expectedName, budget.Name)
	})
	s.Run("ReturnsBudgetWithGivenCurrency", func() {
		expectedCurrency := budgit.GBP
		budget := budgit.NewBudget("", expectedCurrency)
		s.Equal(expectedCurrency, budget.Currency)
	})
}
