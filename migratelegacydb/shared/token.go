package shared

import (
	"crypto/rand"
	"encoding/base64"
	"math"
)

func GenerateToken(l int) string {
	binaryToken := make([]byte, int(math.Ceil(float64(l)*(math.Log(64)/math.Log(255)))))

	_, err := rand.Read(binaryToken)
	if err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(binaryToken)[:l]
}
