package budgit_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

func TestBudgit(t *testing.T) {
	suite.Run(t, new(budgitSuite))
}

type budgitSuite struct {
	suite.Suite
}

func (s *budgitSuite) CMPEqual(expected, actual any, opts ...cmp.Option) {
	if !cmp.Equal(expected, actual, opts...) {
		s.Fail(cmp.Diff(expected, actual, opts...))
	}
}
