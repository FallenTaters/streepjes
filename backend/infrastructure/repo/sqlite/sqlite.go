package sqlite

import (
	"database/sql"
	"sync"

	_ "github.com/glebarez/go-sqlite"
)

func OpenDB(dbname string) (Queryable, error) {
	db, err := sql.Open(`sqlite`, dbname)
	if err != nil {
		return db, err
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		panic(err)
	}

	return &LockedQueryable{Queryable: db}, nil
}

type LockedQueryable struct {
	Queryable
	sync.Mutex
}

func (lq *LockedQueryable) Exec(query string, args ...interface{}) (sql.Result, error) {
	lq.Lock()
	defer lq.Unlock()
	return lq.Queryable.Exec(query, args...)
}

func (lq *LockedQueryable) Query(query string, args ...interface{}) (*sql.Rows, error) {
	lq.Lock()
	defer lq.Unlock()
	return lq.Queryable.Query(query, args...)
}

func (lq *LockedQueryable) QueryRow(query string, args ...interface{}) *sql.Row {
	lq.Lock()
	defer lq.Unlock()
	return lq.Queryable.QueryRow(query, args...)
}

type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}
