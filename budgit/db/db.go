package db

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB struct{}

func New() DB {
	return DB{}
}

type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Execer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

const dbStructKey = "db"

func getAllDBColumns(strct any) []string {
	elem := reflect.TypeOf(strct)
	columns := make([]string, 0, elem.NumField())
	for i := range elem.NumField() {
		columns = append(columns, string(elem.FieldByIndex([]int{i}).Tag.Get(dbStructKey)))
	}
	return columns
}

type idGetter interface {
	GetID() string
}

func mapByID[E idGetter](elems []E) map[string]E {
	elemsByID := make(map[string]E, len(elems))
	for _, elem := range elems {
		elemsByID[elem.GetID()] = elem
	}
	return elemsByID
}
