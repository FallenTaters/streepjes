package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"go.uber.org/zap"
)

func NewUserRepo(db Queryable, logger *zap.Logger) repo.User {
	return &userRepo{
		db:     db,
		logger: logger,
	}
}

type userRepo struct {
	db     Queryable
	logger *zap.Logger
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
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.id = $1;`, id)

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
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.username = $1;`, username) //nolint:lll

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
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.auth_token = $1;`, token) //nolint:lll

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
	if user.Username == `` ||
		user.Name == `` ||
		user.Club == domain.ClubUnknown ||
		user.Role == authdomain.RoleNotAuthorized {
		return repo.ErrUserMissingFields
	}

	res, err := ur.db.Exec(
		`UPDATE users SET username = $1, password = $2, club = $3, name = $4, role = $5, auth_token = $6, auth_time = $7 WHERE id = $8;`,
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

	ur.logger.Info("user updated", zap.Int("id", user.ID), zap.String("username", user.Username))

	return nil
}

func (ur *userRepo) UpdateActivity(user authdomain.User) error {
	res, err := ur.db.Exec(
		`UPDATE users SET auth_token = $1, auth_time = $2 WHERE id = $3;`,
		user.AuthToken, user.AuthTime, user.ID,
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

	if _, ok := ur.getByName(user.Name); ok {
		return 0, fmt.Errorf(`%w, %#v`, repo.ErrUsernameTaken, user)
	}

	if _, ok := ur.GetByUsername(user.Username); ok {
		return 0, fmt.Errorf(`%w, %#v`, repo.ErrUsernameTaken, user)
	}

	row := ur.db.QueryRow(
		`INSERT INTO users (username, password, club, name, role) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		user.Username, user.PasswordHash, user.Club, user.Name, user.Role,
	)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	ur.logger.Info("user created",
		zap.Int("id", id),
		zap.String("username", user.Username),
		zap.String("role", user.Role.String()),
		zap.String("club", user.Club.String()),
	)

	return id, nil
}

func (ur *userRepo) getByName(name string) (authdomain.User, bool) {
	row := ur.db.QueryRow(`SELECT id, username, password, club, name, role, auth_token, auth_time FROM users U WHERE U.name = $1;`, name)

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

func (ur *userRepo) Delete(id int) error {
	res, err := ur.db.Exec(`DELETE FROM users WHERE id = $1;`, id)
	if err != nil {
		if err.Error() == `FOREIGN KEY constraint failed` {
			return repo.ErrUserHasOrders
		}

		panic(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if affected == 0 {
		return repo.ErrUserNotFound
	}

	ur.logger.Info("user deleted", zap.Int("id", id))

	return nil
}
