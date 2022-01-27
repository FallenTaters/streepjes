package static

import "embed"

//go:embed files/*
var assets embed.FS

func Get(name string) ([]byte, error) {
	return assets.ReadFile(`files/` + name)
}
