package backend

import (
	"net/url"
)

func Init(endpoint string) {
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	settings.Endpoint = u
}

type Settings struct {
	Endpoint *url.URL
}

var settings Settings

func (s Settings) URL() string {
	return s.Endpoint.String()
}
