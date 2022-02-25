package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
	"github.com/PotatoesFall/vecty-test/frontend/events"
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

func checkStatus(code int) error {
	if code == http.StatusOK {
		return nil
	}

	if code == http.StatusUnauthorized {
		events.Trigger(events.Unauthorized)
		return ErrUnauthorized
	}

	if code == http.StatusForbidden {
		return ErrForbidden
	}

	return fmt.Errorf(`%w: %d`, ErrStatus, code)
}

func GetCatalog() (api.Catalog, error) {
	resp, err := http.Get(settings.URL() + `/catalog`)
	if err != nil {
		return api.Catalog{}, err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp.StatusCode); err != nil {
		return api.Catalog{}, err
	}

	var catalog api.Catalog
	return catalog, json.NewDecoder(resp.Body).Decode(&catalog)
}

func GetMembers() ([]orderdomain.Member, error) {
	resp, err := http.Get(settings.URL() + `/members`)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp.StatusCode); err != nil {
		return nil, err
	}

	var members []orderdomain.Member
	return members, json.NewDecoder(resp.Body).Decode(&members)
}

func PostLogout() error {
	resp, err := http.Post(settings.URL()+`/logout`, ``, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return fmt.Errorf(`%w: %d`, ErrStatus, resp.StatusCode)
	}

	return nil
}

func PostLogin(req api.Credentials) (authdomain.User, error) {
	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(settings.URL()+`/login`, `application/json`, bytes.NewReader(data))
	if err != nil {
		return authdomain.User{}, err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp.StatusCode); err != nil {
		return authdomain.User{}, err
	}

	var user authdomain.User
	return user, json.NewDecoder(resp.Body).Decode(&user)
}

func PostActive() (authdomain.User, error) {
	resp, err := http.Post(settings.URL()+`/active`, ``, nil)
	if err != nil {
		return authdomain.User{}, err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp.StatusCode); err != nil {
		return authdomain.User{}, err
	}

	var user authdomain.User
	return user, json.NewDecoder(resp.Body).Decode(&user)
}
