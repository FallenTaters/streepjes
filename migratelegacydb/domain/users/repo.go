package users

import (
	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared/buckets"
	"golang.org/x/crypto/bcrypt"
)

func create(user User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = hash

	return buckets.Users.Create(user.Key(), user)
}

func getAll() ([]User, error) {
	users := []User{}

	return users, buckets.Users.GetAll(&User{}, func(ptr interface{}) error {
		users = append(users, *ptr.(*User))
		return nil
	})
}

func get(username string) (User, error) {
	var user User

	return user, buckets.Users.Get([]byte(username), &user)
}

func update(username string, mutate func(User) (User, error)) error {
	return buckets.Users.Update([]byte(username), &User{}, func(ptr interface{}) (interface{}, error) {
		return mutate(*ptr.(*User))
	})
}

func delete(username string) error {
	return buckets.Users.Delete([]byte(username))
}
