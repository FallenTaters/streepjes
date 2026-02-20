package mockdb

import (
	"github.com/FallenTaters/streepjes/domain/authdomain"
)

type User struct {
	GetFunc            func(int) (authdomain.User, error)
	GetAllFunc         func() ([]authdomain.User, error)
	GetByTokenFunc     func(token string) (authdomain.User, error)
	GetByUsernameFunc  func(username string) (authdomain.User, error)
	UpdateFunc         func(user authdomain.User) error
	UpdateActivityFunc func(user authdomain.User) error
	CreateFunc         func(user authdomain.User) (int, error)
	DeleteFunc         func(id int) error
}

func (u User) Get(id int) (authdomain.User, error) {
	return u.GetFunc(id)
}

func (u User) GetAll() ([]authdomain.User, error) {
	return u.GetAllFunc()
}

func (u User) GetByToken(token string) (authdomain.User, error) {
	return u.GetByTokenFunc(token)
}

func (u User) GetByUsername(username string) (authdomain.User, error) {
	return u.GetByUsernameFunc(username)
}

func (u User) Update(user authdomain.User) error {
	return u.UpdateFunc(user)
}

func (u User) UpdateActivity(user authdomain.User) error {
	if u.UpdateActivityFunc != nil {
		return u.UpdateActivityFunc(user)
	}
	return u.UpdateFunc(user)
}

func (u User) Create(user authdomain.User) (int, error) {
	return u.CreateFunc(user)
}

func (u User) Delete(id int) error {
	return u.DeleteFunc(id)
}
