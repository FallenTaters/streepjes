package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(dbname string) (*sql.DB, error) {
	return sql.Open(`sqlite3`, dbname)
}

type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
