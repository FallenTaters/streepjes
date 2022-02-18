package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain"
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

func GetMembers() ([]domain.Member, error) {
	resp, err := http.Get(settings.URL() + `/members`)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp.StatusCode); err != nil {
		return nil, err
	}

	var members []domain.Member
	err = json.NewDecoder(resp.Body).Decode(&members)

	return members, err
}

func PostLogout() error {
	resp, err := http.Post(settings.URL()+`/logout`, ``, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkStatus(resp.StatusCode)
}

func PostLogin(req api.Credentials) error {
	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(settings.URL()+`/login`, `application/json`, bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkStatus(resp.StatusCode)
}
