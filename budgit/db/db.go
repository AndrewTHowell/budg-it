package db

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type DB struct {
	log *zap.SugaredLogger
}

func New(log *zap.SugaredLogger) DB {
	return DB{log: log}
}

type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Execer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func rowsToIDs(rows pgx.Rows) ([]string, error) {
	ids := make([]string, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func structsToPointers[E any](elems []E) []*E {
	ptrElems := make([]*E, 0, len(elems))
	for _, elem := range elems {
		ptrElems = append(ptrElems, &elem)
	}
	return ptrElems
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

type nameGetter interface {
	GetName() string
}

func mapByName[E nameGetter](elems []E) map[string]E {
	elemsByName := make(map[string]E, len(elems))
	for _, elem := range elems {
		elemsByName[elem.GetName()] = elem
	}
	return elemsByName
}
