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
	return &userRepo{db: db, logger: logger}
}

type userRepo struct {
	db     Queryable
	logger *zap.Logger
}

const userColumns = `id, username, password, club, name, role, auth_token, auth_time`

func scanUser(sc interface{ Scan(...any) error }) (authdomain.User, error) {
	var u authdomain.User
	err := sc.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Club, &u.Name, &u.Role, &u.AuthToken, &u.AuthTime)
	return u, err
}

func (ur *userRepo) GetAll() ([]authdomain.User, error) {
	rows, err := ur.db.Query(`SELECT ` + userColumns + ` FROM users;`)
	if err != nil {
		return nil, fmt.Errorf("userRepo.GetAll: query: %w", err)
	}
	defer rows.Close()

	var users []authdomain.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetAll: scan: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("userRepo.GetAll: rows: %w", err)
	}

	return users, nil
}

func (ur *userRepo) Get(id int) (authdomain.User, error) {
	u, err := scanUser(ur.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE id = $1;`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return authdomain.User{}, repo.ErrUserNotFound
	}
	if err != nil {
		return authdomain.User{}, fmt.Errorf("userRepo.Get(%d): %w", id, err)
	}
	return u, nil
}

func (ur *userRepo) GetByUsername(username string) (authdomain.User, error) {
	u, err := scanUser(ur.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE username = $1;`, username))
	if errors.Is(err, sql.ErrNoRows) {
		return authdomain.User{}, repo.ErrUserNotFound
	}
	if err != nil {
		return authdomain.User{}, fmt.Errorf("userRepo.GetByUsername: %w", err)
	}
	return u, nil
}

func (ur *userRepo) GetByToken(token string) (authdomain.User, error) {
	u, err := scanUser(ur.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE auth_token = $1;`, token))
	if errors.Is(err, sql.ErrNoRows) {
		return authdomain.User{}, repo.ErrUserNotFound
	}
	if err != nil {
		return authdomain.User{}, fmt.Errorf("userRepo.GetByToken: %w", err)
	}
	return u, nil
}

func (ur *userRepo) Update(user authdomain.User) error {
	if user.Username == `` || user.Name == `` || user.Club == domain.ClubUnknown || user.Role == authdomain.RoleNotAuthorized {
		return repo.ErrUserMissingFields
	}

	res, err := ur.db.Exec(
		`UPDATE users SET username = $1, password = $2, club = $3, name = $4, role = $5, auth_token = $6, auth_time = $7 WHERE id = $8;`,
		user.Username, user.PasswordHash, user.Club, user.Name, user.Role, user.AuthToken, user.AuthTime, user.ID,
	)
	if err != nil {
		return fmt.Errorf("userRepo.Update(%d): %w", user.ID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("userRepo.Update(%d): rows affected: %w", user.ID, err)
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
		return fmt.Errorf("userRepo.UpdateActivity(%d): %w", user.ID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("userRepo.UpdateActivity(%d): rows affected: %w", user.ID, err)
	}
	if affected == 0 {
		return repo.ErrUserNotFound
	}

	return nil
}

func (ur *userRepo) Create(user authdomain.User) (int, error) {
	if user.Username == `` || len(user.PasswordHash) == 0 || user.Club == domain.ClubUnknown || user.Name == `` || user.Role == authdomain.RoleNotAuthorized {
		return 0, fmt.Errorf(`%w: %#v`, repo.ErrUserMissingFields, user)
	}

	_, err := ur.getByName(user.Name)
	if err == nil {
		return 0, fmt.Errorf(`%w: name %q`, repo.ErrUsernameTaken, user.Name)
	}
	if !errors.Is(err, repo.ErrUserNotFound) {
		return 0, fmt.Errorf("userRepo.Create: check name: %w", err)
	}

	_, err = ur.GetByUsername(user.Username)
	if err == nil {
		return 0, fmt.Errorf(`%w: username %q`, repo.ErrUsernameTaken, user.Username)
	}
	if !errors.Is(err, repo.ErrUserNotFound) {
		return 0, fmt.Errorf("userRepo.Create: check username: %w", err)
	}

	var id int
	if err := ur.db.QueryRow(
		`INSERT INTO users (username, password, club, name, role) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		user.Username, user.PasswordHash, user.Club, user.Name, user.Role,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("userRepo.Create: insert: %w", err)
	}

	ur.logger.Info("user created",
		zap.Int("id", id),
		zap.String("username", user.Username),
		zap.String("role", user.Role.String()),
		zap.String("club", user.Club.String()),
	)

	return id, nil
}

func (ur *userRepo) getByName(name string) (authdomain.User, error) {
	u, err := scanUser(ur.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE name = $1;`, name))
	if errors.Is(err, sql.ErrNoRows) {
		return authdomain.User{}, repo.ErrUserNotFound
	}
	if err != nil {
		return authdomain.User{}, fmt.Errorf("userRepo.getByName: %w", err)
	}
	return u, nil
}

func (ur *userRepo) Delete(id int) error {
	res, err := ur.db.Exec(`DELETE FROM users WHERE id = $1;`, id)
	if err != nil {
		if err.Error() == `FOREIGN KEY constraint failed` {
			return repo.ErrUserHasOrders
		}
		return fmt.Errorf("userRepo.Delete(%d): %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("userRepo.Delete(%d): rows affected: %w", id, err)
	}
	if affected == 0 {
		return repo.ErrUserNotFound
	}

	ur.logger.Info("user deleted", zap.Int("id", id))
	return nil
}
