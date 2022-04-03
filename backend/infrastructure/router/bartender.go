package router

import (
	"net/http"
	"strconv"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/labstack/echo/v4"
)

func bartenderRoutes(r *echo.Group, orderService order.Service) {
	r.GET(`/catalog`, getCatalog)
	r.GET(`/members`, getMembers(orderService))
	r.GET(`/member/:id`, getMember(orderService))
	r.POST(`/order`, postOrder(orderService))
}

func getCatalog(c echo.Context) error { //nolint:funlen
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

func getMembers(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, orderService.GetAllMembers())
	}
}

func getMember(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param(`id`))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		member, ok := orderService.GetMemberDetails(id)
		if !ok {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, member)
	}
}

func postOrder(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		order, ok := readJSON[orderdomain.Order](c)
		if !ok {
			return nil
		}

		if err := orderService.PlaceOrder(order, userFromContext(c)); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}
