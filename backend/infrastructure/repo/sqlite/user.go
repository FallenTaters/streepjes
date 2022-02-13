package sqlite

import (
	"database/sql"
	"errors"

	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain"
)

func NewUserRepo(db *sql.DB) repo.User {
	return &userRepo{
		User: nil,
		db:   db,
	}
}

type userRepo struct {
	repo.User // TODO remove
	db        *sql.DB
}

func (ur *userRepo) GetByUsername(username string) (domain.User, bool) {
	row := ur.db.QueryRow(`SELECT * FROM users U WHERE U.username = ?;`, username)

	err := row.Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, false
	}
	if err != nil {
		panic(err)
	}

	return domain.User{}, true // TODO
}
