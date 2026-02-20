package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/authdomain"
)

var (
	ErrPasswordEmpty      = errors.New("new password cannot be empty")
	ErrPasswordWrong      = errors.New("original password is incorrect")
	ErrNameEmpty          = errors.New("name cannot be empty")
	ErrUserHasOrders      = errors.New("cannot delete user with orders")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type Service interface {
	Login(user, pass string) (authdomain.User, error)
	Check(token string) (authdomain.User, error)
	Active(id int) error
	Logout(id int) error
	Register(user authdomain.User, password string) error
	Update(user authdomain.User, password string) error
	ChangePassword(user authdomain.User, changePassword api.ChangePassword) error
	ChangeName(user authdomain.User, name string) error
	GetUsers() ([]authdomain.User, error)
	Delete(id int) error
}

func New(userRepo repo.User, orderRepo repo.Order) Service {
	return &service{userRepo, orderRepo}
}

type service struct {
	users  repo.User
	orders repo.Order
}

func (s *service) Login(username, pass string) (authdomain.User, error) {
	user, err := s.users.GetByUsername(username)
	if errors.Is(err, repo.ErrUserNotFound) {
		return authdomain.User{}, ErrInvalidCredentials
	}
	if err != nil {
		return authdomain.User{}, fmt.Errorf("auth.Login: %w", err)
	}

	if !CheckPassword(user.PasswordHash, pass) {
		return authdomain.User{}, ErrInvalidCredentials
	}

	user.AuthToken = generateToken()
	user.AuthTime = time.Now()

	if err := s.users.UpdateActivity(user); err != nil {
		return authdomain.User{}, fmt.Errorf("auth.Login: update activity: %w", err)
	}

	return user, nil
}

func (s *service) Check(token string) (authdomain.User, error) {
	if token == `` {
		return authdomain.User{}, ErrInvalidToken
	}

	user, err := s.users.GetByToken(token)
	if errors.Is(err, repo.ErrUserNotFound) {
		return authdomain.User{}, ErrInvalidToken
	}
	if err != nil {
		return authdomain.User{}, fmt.Errorf("auth.Check: %w", err)
	}

	if time.Since(user.AuthTime) > authdomain.TokenDuration {
		return authdomain.User{}, ErrInvalidToken
	}

	user.AuthTime = time.Now()
	if err := s.users.UpdateActivity(user); err != nil {
		return authdomain.User{}, fmt.Errorf("auth.Check: update activity: %w", err)
	}

	return user, nil
}

func (s *service) Active(id int) error {
	user, err := s.users.Get(id)
	if errors.Is(err, repo.ErrUserNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("auth.Active: %w", err)
	}

	user.AuthTime = time.Now()

	if err := s.users.UpdateActivity(user); err != nil {
		return fmt.Errorf("auth.Active: update: %w", err)
	}
	return nil
}

func (s *service) Logout(id int) error {
	user, err := s.users.Get(id)
	if errors.Is(err, repo.ErrUserNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("auth.Logout: %w", err)
	}

	user.AuthToken = ``
	user.AuthTime = time.Now().Add(-authdomain.TokenDuration)

	if err := s.users.UpdateActivity(user); err != nil {
		return fmt.Errorf("auth.Logout: update: %w", err)
	}
	return nil
}

func (s *service) Register(user authdomain.User, password string) error {
	user.PasswordHash = HashPassword(password)

	_, err := s.users.Create(user)
	return err
}

func (s *service) Update(userChanges authdomain.User, password string) error {
	user, err := s.users.Get(userChanges.ID)
	if err != nil {
		return fmt.Errorf("auth.Update: %w", err)
	}

	if password != `` {
		user.PasswordHash = HashPassword(password)
	}

	user.Name = userChanges.Name
	user.Username = userChanges.Username
	user.Club = userChanges.Club
	user.Role = userChanges.Role

	return s.users.Update(user)
}

func (s *service) ChangePassword(user authdomain.User, changePassword api.ChangePassword) error {
	if changePassword.New == `` {
		return ErrPasswordEmpty
	}

	if !CheckPassword(user.PasswordHash, changePassword.Original) {
		return ErrPasswordWrong
	}

	user.PasswordHash = HashPassword(changePassword.New)

	return s.users.Update(user)
}

func (s *service) ChangeName(user authdomain.User, name string) error {
	if name == `` {
		return ErrNameEmpty
	}

	user.Name = name

	return s.users.Update(user)
}

func (s *service) GetUsers() ([]authdomain.User, error) {
	return s.users.GetAll()
}

func (s *service) Delete(id int) error {
	orders, err := s.orders.Filter(repo.OrderFilter{
		BartenderID: id,
	})
	if err != nil {
		return fmt.Errorf("auth.Delete: check orders: %w", err)
	}

	if len(orders) > 0 {
		return ErrUserHasOrders
	}

	return s.users.Delete(id)
}
