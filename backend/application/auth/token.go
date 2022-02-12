package auth

import (
	"crypto/rand"
	"encoding/base64"
)

const tokenLen = 32

func generateToken() string {
	binaryToken := make([]byte, base64.RawURLEncoding.DecodedLen(tokenLen)+1)

	_, err := rand.Read(binaryToken)
	if err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(binaryToken)[:tokenLen]
}
