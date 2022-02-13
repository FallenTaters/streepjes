package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func checkPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
