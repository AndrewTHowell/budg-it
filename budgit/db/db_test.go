package db_test

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

func TestDB(t *testing.T) {
	suite.Run(t, new(dbSuite))
}

type dbSuite struct {
	suite.Suite

	pgContainer testcontainers.Container

	db   db.DB
	conn *pgx.Conn
}

func (s *dbSuite) SetupSuite() {
	migrationFilePaths := []string{}
	migrationsPath := "../migrations"
	err := filepath.WalkDir(migrationsPath, func(path string, d fs.DirEntry, _ error) error {
		if strings.HasSuffix(path, ".up.sql") {
			migrationFilePaths = append(migrationFilePaths, path)
		}
		return nil
	})
	s.Require().NoError(err, "unexpected error getting migration file paths")

	pgContainer, err := postgres.Run(context.Background(), "postgres:16-alpine",
		postgres.WithInitScripts(migrationFilePaths...),
		postgres.WithDatabase("budgit"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	s.Require().NoError(err, "unexpected error running postgres container")
	s.pgContainer = pgContainer

	connString, err := pgContainer.ConnectionString(context.Background(), "sslmode=disable")
	s.Require().NoError(err, "unexpected error getting postgres container URL")
	conn, err := pgx.Connect(context.Background(), connString)
	s.Require().NoError(err, "unexpected error connecting to postgres container")
	s.conn = conn

	s.db = db.New(zap.NewNop().Sugar())
}

func (s *dbSuite) TearDownTest() {
	s.truncateTables("accounts", "payees", "transactions")
}

func (s *dbSuite) TearDownSuite() {
	s.Require().NoError(s.conn.Close(context.Background()), "unexpected error closing connection")
	s.Require().NoError(s.pgContainer.Terminate(context.Background()), "unexpected error terminating postgres container")
}

func (s *dbSuite) truncateTables(tables ...string) {
	_, err := s.conn.Exec(context.Background(), fmt.Sprintf(`TRUNCATE TABLE %s`, strings.Join(tables, ", ")))
	s.Require().NoError(err, "unexpected error truncating tables %+v", tables)
}

func (s *dbSuite) CMPEqual(expected, actual any, opts ...cmp.Option) {
	if !cmp.Equal(expected, actual, opts...) {
		s.Fail(cmp.Diff(expected, actual, opts...))
	}
}
