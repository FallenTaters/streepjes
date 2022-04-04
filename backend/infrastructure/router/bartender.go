package router

import (
	"net/http"
	"strconv"

	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/labstack/echo/v4"
)

func bartenderRoutes(r *echo.Group, orderService order.Service) {
	r.GET(`/catalog`, getCatalog(orderService))
	r.GET(`/members`, getMembers(orderService))
	r.GET(`/member/:id`, getMember(orderService))
	r.POST(`/order`, postOrder(orderService))
}

func getCatalog(orderService order.Service) func(echo.Context) error {
	return func(c echo.Context) error {
		catalog := orderService.GetCatalog()
		return c.JSON(http.StatusOK, catalog)
	}
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
