package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
)

func NewUserRepo(db *sql.DB) repo.User {
	return &userRepo{
		db: db,
	}
}

type userRepo struct {
	db *sql.DB
}

func (ur *userRepo) GetAll() []authdomain.User {
	rows, err := ur.db.Query(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var user authdomain.User
	users := make([]authdomain.User, 0)

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Club, &user.Name, &user.Role, &user.AuthToken, &user.AuthTime)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return users
}

func (ur *userRepo) Get(id int) (authdomain.User, bool) {
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.id = ?;`, id)

	var user authdomain.User

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Club, &user.Name, &user.Role, &user.AuthToken, &user.AuthTime)
	if errors.Is(err, sql.ErrNoRows) {
		return authdomain.User{}, false
	}
	if err != nil {
		panic(err)
	}

	return user, true
}

func (ur *userRepo) GetByUsername(username string) (authdomain.User, bool) {
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.username = ?;`, username) //nolint:lll

	var user authdomain.User

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Club, &user.Name, &user.Role, &user.AuthToken, &user.AuthTime)
	if errors.Is(err, sql.ErrNoRows) {
		return authdomain.User{}, false
	}
	if err != nil {
		panic(err)
	}

	return user, true
}

func (ur *userRepo) GetByToken(token string) (authdomain.User, bool) {
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.auth_token = ?;`, token) //nolint:lll

	var user authdomain.User

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Club, &user.Name, &user.Role, &user.AuthToken, &user.AuthTime)
	if errors.Is(err, sql.ErrNoRows) {
		return authdomain.User{}, false
	}
	if err != nil {
		panic(err)
	}

	return user, true
}

func (ur *userRepo) Update(user authdomain.User) error {
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

func (ur *userRepo) Create(user authdomain.User) (int, error) {
	if user.Username == `` ||
		len(user.PasswordHash) == 0 ||
		user.Club == domain.ClubUnknown ||
		user.Name == `` ||
		user.Role == authdomain.RoleNotAuthorized {
		return 0, fmt.Errorf(`%w: %#v`, repo.ErrUserMissingFields, user)
	}

	res, err := ur.db.Exec(
		`INSERT INTO users (username, password, club, name, role) VALUES (?, ?, ?, ?, ?);`,
		user.Username, user.PasswordHash, user.Club, user.Name, user.Role, user.AuthToken, user.AuthTime, user.ID,
	)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return int(id), nil
}

func (ur *userRepo) Delete(id int) error {
	res, err := ur.db.Exec(`DELETE FROM users WHERE id = ?;`, id)
	if err != nil {
		panic(err) // TODO, is this where we panic when foreignkey conflict? or do we just get 0 rows affected?
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
