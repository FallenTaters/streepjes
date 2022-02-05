package router

import (
	"net/http"
	"strings"

	"git.fuyu.moe/Fuyu/router"
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain"
)

type Static func(filename string) ([]byte, error)

func New(static Static) *router.Router {
	r := router.New()

	r.ErrorHandler = panicHandler

	r.GET(`/version`, getVersion)
	r.GET(`/`, getIndex(static))
	r.GET(`/static/*name`, getStatic(static))

	r.GET(`/catalog`, getCatalog)

	return r
}

func getVersion(c *router.Context) error {
	return c.String(http.StatusOK, version())
}

func getIndex(assets Static) router.Handle {
	return func(c *router.Context) error {
		index, err := assets(`index.html`)
		if err != nil {
			panic(err)
		}

		c.Response.Header().Set(`Content-Type`, `text/html`)
		return c.Bytes(http.StatusOK, index)
	}
}

func getStatic(assets Static) router.Handle {
	return func(c *router.Context) error {
		name := strings.TrimPrefix(c.Param(`name`), `/`)
		asset, err := assets(name)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		return c.Bytes(http.StatusOK, asset)
	}
}

func getCatalog(c *router.Context) error {
	// TODO make actual catalog
	catalog := api.Catalog{
		Categories: []domain.Category{
			{
				ID:   1,
				Name: `Food`,
			},
			{
				ID:   2,
				Name: `Drinks`,
			},
		},
		Items: []domain.Item{
			{
				ID:              1,
				CategoryID:      1,
				Name:            `Chicken Fingers`,
				PriceGladiators: 250,
				PriceParabool:   200,
			},
			{
				ID:              2,
				CategoryID:      1,
				Name:            `Snickers`,
				PriceGladiators: 200,
				PriceParabool:   150,
			},

			{
				ID:              3,
				CategoryID:      2,
				Name:            `Beer`,
				PriceGladiators: 150,
				PriceParabool:   120,
			},
			{
				ID:              4,
				CategoryID:      2,
				Name:            `Wine`,
				PriceGladiators: 250,
				PriceParabool:   220,
			},
			{
				ID:              5,
				CategoryID:      1,
				Name:            `Chicken Fingers 2`,
				PriceGladiators: 250,
				PriceParabool:   200,
			},
			{
				ID:              6,
				CategoryID:      1,
				Name:            `Snickers 2`,
				PriceGladiators: 200,
				PriceParabool:   150,
			},
			{
				ID:              7,
				CategoryID:      1,
				Name:            `Chicken Fingers 3`,
				PriceGladiators: 250,
				PriceParabool:   200,
			},
			{
				ID:              8,
				CategoryID:      1,
				Name:            `Snickers 3`,
				PriceGladiators: 200,
				PriceParabool:   150,
			},
		},
	}
	return c.JSON(http.StatusOK, catalog)
}
