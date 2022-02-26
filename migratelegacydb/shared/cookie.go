package shared

import (
	"encoding/base64"
	"net/http"
)

var settings = struct {
	DisableSecure bool
}{
	DisableSecure: false,
}

func Init(disableSecure bool) {
	settings.DisableSecure = disableSecure
}

func GetCookie(req *http.Request, name string) ([]byte, bool) {
	cookie, err := req.Cookie(name)
	if err != nil {
		return []byte{}, false
	}

	v, err := base64.RawStdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return []byte{}, false
	}

	return v, true
}

func SetCookie(w http.ResponseWriter, name string, value []byte, maxAgeInSeconds int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    base64.RawStdEncoding.EncodeToString(value),
		HttpOnly: true,
		Path:     `/`,
		Secure:   !settings.DisableSecure,
		MaxAge:   maxAgeInSeconds,
	})
}

func UnsetCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    ``,
		HttpOnly: true,
		Secure:   !settings.DisableSecure,
		MaxAge:   -1,
	})
}
