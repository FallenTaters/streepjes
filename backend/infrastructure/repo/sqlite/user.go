package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain"
)

func NewUserRepo(db *sql.DB) repo.User {
	return &userRepo{
		db: db,
	}
}

type userRepo struct {
	db *sql.DB
}

func (ur *userRepo) Get(id int) (domain.User, bool) {
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.id = ?;`, id)

	var user domain.User

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Club, &user.Name, &user.Role, &user.AuthToken, &user.AuthTime)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, false
	}
	if err != nil {
		panic(err)
	}

	return user, true
}

func (ur *userRepo) GetByUsername(username string) (domain.User, bool) {
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.username = ?;`, username) //nolint:lll

	var user domain.User

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Club, &user.Name, &user.Role, &user.AuthToken, &user.AuthTime)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, false
	}
	if err != nil {
		panic(err)
	}

	return user, true
}

func (ur *userRepo) GetByToken(token string) (domain.User, bool) {
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.auth_token = ?;`, token) //nolint:lll

	var user domain.User

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Club, &user.Name, &user.Role, &user.AuthToken, &user.AuthTime)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, false
	}
	if err != nil {
		panic(err)
	}

	return user, true
}

func (ur *userRepo) Update(user domain.User) error {
	res, err := ur.db.Exec(
		`UPDATE users SET username = ?, password = ?, club = ?, name = ?, role = ?, auth_token = ?, auth_time = ? WHERE id = ?;`,
		user.Username, user.PasswordHash, user.Club, user.Name, user.Role, user.AuthToken, user.AuthTime, user.ID,
	)
	if err != nil {
		panic(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if affected == 0 {
		return repo.ErrUserNotFound
	}

	return nil
}

func (ur *userRepo) Create(user domain.User) error {
	if user.Username == `` || len(user.PasswordHash) == 0 || user.Club == domain.ClubUnknown || user.Name == `` || user.Role == domain.RoleNotAuthorized {
		return fmt.Errorf(`User not filled for creating: %#v`, user)
	}

	res, err := ur.db.Exec(
		`INSERT INTO users (username, password, club, name, role) VALUES (?, ?, ?, ?, ?);`,
		user.Username, user.PasswordHash, user.Club, user.Name, user.Role, user.AuthToken, user.AuthTime, user.ID,
	)
	if err != nil {
		panic(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if affected == 0 {
		return repo.ErrUserNotFound
	}

	return nil
}
