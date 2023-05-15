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
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}
