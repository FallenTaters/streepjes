package backend

import (
	"encoding/json"
	"net/http"

	"github.com/PotatoesFall/vecty-test/api"
)

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
