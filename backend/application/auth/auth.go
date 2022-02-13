package auth

import (
	"time"

	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain"
)

const tokenDuration = 5 * time.Minute

type Service interface {
	// Login returns the user if credentials are correct
	// otherwise, it returns false
	Login(user, pass string) (domain.User, bool)

	// Check gets the user with the correct token
	// if the token is expired or unknown, it returns false
	Check(token string) (domain.User, bool)

	// Active refreshes a users token, setting AuthTime to now
	// if the user is not found, it is a no-op
	Active(id int)

	// Logout deletes the users AuthToken, if found
	// if the user is not found, it is a no-op
	Logout(id int)
}

func New(userRepo repo.User) Service {
	return &service{userRepo}
}

type service struct {
	users repo.User
}

func (s *service) Login(username, pass string) (domain.User, bool) {
	user, ok := s.users.GetByUsername(username)
	if !ok {
		return domain.User{}, false
	}

	if !checkPassword(user.Password, pass) {
		return domain.User{}, false
	}

	user.AuthToken = generateToken()

	err := s.users.Update(user) //nolint:ifshort
	if err != nil {
		panic(err)
	}

	return user, true
}

func (s *service) Check(token string) (domain.User, bool) {
	if token == `` {
		return domain.User{}, false
	}

	user, ok := s.users.GetByToken(token)
	if !ok {
		return domain.User{}, false
	}

	if time.Since(user.AuthTime) > tokenDuration {
		return domain.User{}, false
	}

	return user, true
}

func (s *service) Active(id int) {
	user, ok := s.users.Get(id)
	if !ok {
		return
	}

	user.AuthTime = time.Now()

	_ = s.users.Update(user)
}

func (s *service) Logout(id int) {
	user, ok := s.users.Get(id)
	if !ok {
		return
	}

	user.AuthToken = ``
	user.AuthTime = time.Now().Add(-tokenDuration)

	_ = s.users.Update(user)
}
