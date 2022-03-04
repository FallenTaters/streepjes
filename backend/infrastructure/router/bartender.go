package router

import (
	"net/http"

	"git.fuyu.moe/Fuyu/router"
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func bartenderRoutes(r *router.Group, orderService order.Service) {
	r.GET(`/catalog`, getCatalog)
	r.GET(`/members`, getMembers(orderService))
}

func getCatalog(c *router.Context) error { //nolint:funlen
	// TODO make actual catalog
	catalog := api.Catalog{
		Categories: []orderdomain.Category{
			{
				ID:   1,
				Name: `Food`,
			},
			{
				ID:   2,
				Name: `Drinks`,
			},
			{
				ID:   3,
				Name: `Empty`,
			},
			{
				ID:   4,
				Name: `Empty for one club`,
			},
		},
		Items: []orderdomain.Item{
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
				PriceParabool:   0,
			},
			{
				ID:              9,
				CategoryID:      4,
				Name:            `Snickers 4`,
				PriceGladiators: 200,
				PriceParabool:   0,
			},
		},
	}

	return c.JSON(http.StatusOK, catalog)
}

func getMembers(orderService order.Service) router.Handle {
	return func(c *router.Context) error {
		members := orderService.GetAllMembers()

		return c.JSON(http.StatusOK, members)
	}
}
