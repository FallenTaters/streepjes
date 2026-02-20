package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func OpenDB(connectionString string) (Queryable, error) {
	db, err := sql.Open(`postgres`, connectionString)
	if err != nil {
		return db, err
	}

	return db, nil
}

type Queryable interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Close() error
}
