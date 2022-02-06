package backend

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain"
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
	resp, err := http.Get(settings.URL() + `/catalog`)
	if err != nil {
		return api.Catalog{}, err
	}
	defer resp.Body.Close()

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

	var members []domain.Member
	err = json.NewDecoder(resp.Body).Decode(&members)
	return members, err
}
