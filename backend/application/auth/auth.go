package auth

import (
	"time"

	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
)

type Service interface {
	// Login returns the user if credentials are correct
	// otherwise, it returns false
	Login(user, pass string) (authdomain.User, bool)

	// Check gets the user with the correct token
	// if the token is expired or unknown, it returns false
	Check(token string) (authdomain.User, bool)

	// Active refreshes a users token, setting AuthTime to now
	// if the user is not found, it is a no-op
	Active(id int)

	// Logout deletes the users AuthToken, if found
	// if the user is not found, it is a no-op
	Logout(id int)

	// Register registers a new user. It sets the passwordHash and ID
	// if the username is taken, it return repo.ErrUsernameTaken
	Register(user authdomain.User, password string) error
}

func New(userRepo repo.User) Service {
	return &service{userRepo}
}

type service struct {
	users repo.User
}

func (s *service) Login(username, pass string) (authdomain.User, bool) {
	user, ok := s.users.GetByUsername(username)
	if !ok {
		return authdomain.User{}, false
	}

	if !checkPassword(user.PasswordHash, pass) {
		return authdomain.User{}, false
	}

	user.AuthToken = generateToken()
	user.AuthTime = time.Now()

	err := s.users.Update(user) //nolint:ifshort
	if err != nil {
		panic(err)
	}

	return user, true
}

func (s *service) Check(token string) (authdomain.User, bool) {
	if token == `` {
		return authdomain.User{}, false
	}

	user, ok := s.users.GetByToken(token)
	if !ok {
		return authdomain.User{}, false
	}

	if time.Since(user.AuthTime) > authdomain.TokenDuration {
		return authdomain.User{}, false
	}

	user.AuthTime = time.Now()
	err := s.users.Update(user)
	if err != nil {
		return authdomain.User{}, false
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
	user.AuthTime = time.Now().Add(-authdomain.TokenDuration)

	_ = s.users.Update(user)
}

func (s *service) Register(user authdomain.User, password string) error {
	user.PasswordHash = hashPassword(password)

	return s.users.Create(user)
}
