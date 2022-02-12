package auth

import "golang.org/x/crypto/bcrypt"

func checkPassword(hash []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) != nil
}
