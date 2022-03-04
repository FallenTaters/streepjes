package backend

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

var (
	ErrStatus       = errors.New(`received unexpected status code`)
	ErrUnauthorized = errors.New(`request not authorized`)
	ErrForbidden    = errors.New(`request forbidden`)
)

func Init(endpoint *url.URL) {
	settings.Endpoint = endpoint
}

type Settings struct {
	Endpoint *url.URL
}

var settings Settings

func (s Settings) URL() string {
	return s.Endpoint.String()
}

func GetCatalog() (api.Catalog, error) {
	var catalog api.Catalog
	return catalog, get(`/catalog`, &catalog)
}

func GetMembers() ([]orderdomain.Member, error) {
	var members []orderdomain.Member
	return members, get(`/members`, &members)
}

func PostLogout() error {
	resp, err := http.Post(settings.URL()+`/logout`, ``, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// don't check status because logout is often called when already logged out, just to ensure logout

	return nil
}

func PostLogin(req api.Credentials) (authdomain.User, error) {
	var user authdomain.User
	return user, post(`/login`, req, &user)
}

func PostActive() (authdomain.User, error) {
	var user authdomain.User
	return user, post(`/active`, nil, &user)
}

func PostChangePassword(changePassword api.ChangePassword) error {
	return post(`/me/password`, changePassword, nil)
}

func PostChangeName(name string) error {
	return post(`/me/name`, name, nil)
}
