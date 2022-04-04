package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(dbname string) (*sql.DB, error) {
	db, err := sql.Open(`sqlite3`, dbname)
	if err != nil {
		return db, err
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		panic(err)
	}

	return db, nil
}

type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
