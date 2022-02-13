package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain"
)

var ErrStatus = errors.New(`received unexpected status code`)

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
	resp, err := http.Get(settings.URL() + `/catalog`)
	if err != nil {
		return api.Catalog{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return api.Catalog{}, fmt.Errorf(`%w: %d`, ErrStatus, resp.StatusCode)
	}

	var catalog api.Catalog
	err = json.NewDecoder(resp.Body).Decode(&catalog)

	return catalog, err
}

func GetMembers() ([]domain.Member, error) {
	resp, err := http.Get(settings.URL() + `/members`)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`%w: %d`, ErrStatus, resp.StatusCode)
	}

	var members []domain.Member
	err = json.NewDecoder(resp.Body).Decode(&members)

	return members, err
}

func Logout() error {
	resp, err := http.Post(settings.URL()+`/logout`, ``, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf(`%w: %d`, ErrStatus, resp.StatusCode)
	}

	return nil
}
