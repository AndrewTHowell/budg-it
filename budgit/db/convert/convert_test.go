package convert_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

func TestConvert(t *testing.T) {
	suite.Run(t, new(convertSuite))
}

type convertSuite struct {
	suite.Suite

	pgContainer testcontainers.Container
	conn        *pgx.Conn
}

func (s *convertSuite) CMPEqual(expected, actual any, opts ...cmp.Option) {
	if !cmp.Equal(expected, actual, opts...) {
		s.Fail(cmp.Diff(expected, actual, opts...))
	}
}
