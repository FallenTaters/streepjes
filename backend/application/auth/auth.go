package auth

import (
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain"
)

type Service interface {
	Login(user, pass string) (domain.User, bool)

	Check(token string) domain.User
	Active(domain.User)

	Logout(domain.User)
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

	err := s.users.Update(user)
	if err != nil {
		panic(err)
	}

	return user, true
}
