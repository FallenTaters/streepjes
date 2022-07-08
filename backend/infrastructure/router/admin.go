package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/labstack/echo/v4"
)

func adminRoutes(r *echo.Group, authService auth.Service, orderService order.Service) {
	r.GET(`/users`, getUsers(authService))
	r.POST(`/users/new`, postNewUser(authService))
	r.POST(`/users/edit`, postEditUser(authService))
	r.POST(`/users/:id/delete`, postDeleteUser(authService))

	r.POST(`/category/new`, postNewCategory(orderService))
	r.POST(`/category/update`, postUpdateCategory(orderService))
	r.POST(`/category/:id/delete`, postDeleteCategory(orderService))

	r.POST(`/item/new`, postNewItem(orderService))
	r.POST(`/item/update`, postUpdateItem(orderService))
	r.POST(`/item/:id/delete`, postDeleteItem(orderService))

	r.POST(`/members/new`, postNewMember(orderService))
	r.POST(`/members/edit`, postEditMember(orderService))
	r.POST(`/members/:id/delete`, postDeleteMember(orderService))

	r.GET(`/billing/orders`, getBillingOrders(orderService))
	r.GET(`/download`, getDownload(orderService))
}

func getUsers(authService auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		users := authService.GetUsers()
		return c.JSON(http.StatusOK, users)
	}
}

func postNewUser(authService auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := readJSON[api.UserWithPassword](c)
		if !ok {
			return nil
		}

		err := authService.Register(user.User, user.Password)
		if err != nil {
			if errors.Is(err, repo.ErrUsernameTaken) {
				return c.String(http.StatusBadRequest, repo.ErrUsernameTaken.Error())
			}

			if errors.Is(err, repo.ErrUserMissingFields) {
				return c.String(http.StatusBadRequest, repo.ErrUserMissingFields.Error())
			}

			return c.NoContent(http.StatusBadRequest)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postEditUser(authService auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := readJSON[api.UserWithPassword](c)
		if !ok {
			return nil
		}

		fmt.Printf("%#v\n", user)
		err := authService.Update(user.User, user.Password)
		if err != nil {
			if errors.Is(err, repo.ErrUserNotFound) {
				return c.String(http.StatusBadRequest, repo.ErrUserNotFound.Error())
			}

			if errors.Is(err, repo.ErrUsernameTaken) {
				return c.String(http.StatusBadRequest, repo.ErrUsernameTaken.Error())
			}

			if errors.Is(err, repo.ErrUserMissingFields) {
				return c.String(http.StatusBadRequest, repo.ErrUserMissingFields.Error())
			}

			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postDeleteUser(authService auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param(`id`))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		ok := authService.Delete(id)
		if !ok {
			return c.NoContent(http.StatusBadRequest)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postNewCategory(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		cat, ok := readJSON[orderdomain.Category](c)
		if !ok {
			return nil
		}

		if err := orderService.NewCategory(cat); err != nil {
			if errors.Is(err, repo.ErrCategoryNameTaken) {
				return c.String(http.StatusBadRequest, repo.ErrCategoryNameTaken.Error())
			}

			if errors.Is(err, repo.ErrCategoryNameEmpty) {
				return c.String(http.StatusBadRequest, repo.ErrCategoryNameEmpty.Error())
			}

			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postUpdateCategory(orderService order.Service) echo.HandlerFunc { //nolint:dupl
	return func(c echo.Context) error {
		cat, ok := readJSON[orderdomain.Category](c)
		if !ok {
			return nil
		}

		if err := orderService.UpdateCategory(cat); err != nil {
			if errors.Is(err, repo.ErrCategoryNameTaken) {
				return c.String(http.StatusBadRequest, repo.ErrCategoryNameTaken.Error())
			}

			if errors.Is(err, repo.ErrCategoryNameEmpty) {
				return c.String(http.StatusBadRequest, repo.ErrCategoryNameEmpty.Error())
			}

			if errors.Is(err, repo.ErrCategoryNotFound) {
				return c.String(http.StatusBadRequest, repo.ErrCategoryNotFound.Error())
			}

			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postDeleteCategory(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param(`id`))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		if err := orderService.DeleteCategory(id); err != nil {
			if errors.Is(err, repo.ErrCategoryNotFound) {
				return c.String(http.StatusBadRequest, repo.ErrCategoryNotFound.Error())
			}

			if errors.Is(err, repo.ErrCategoryHasItems) {
				return c.String(http.StatusBadRequest, repo.ErrCategoryHasItems.Error())
			}

			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postNewItem(orderService order.Service) echo.HandlerFunc { //nolint:dupl
	return func(c echo.Context) error {
		item, ok := readJSON[orderdomain.Item](c)
		if !ok {
			return nil
		}

		if err := orderService.NewItem(item); err != nil {
			if errors.Is(err, repo.ErrItemNameTaken) {
				return c.String(http.StatusBadRequest, repo.ErrItemNameTaken.Error())
			}

			if errors.Is(err, repo.ErrItemNameEmpty) {
				return c.String(http.StatusBadRequest, repo.ErrItemNameEmpty.Error())
			}

			if errors.Is(err, repo.ErrCategoryNotFound) {
				return c.String(http.StatusBadRequest, repo.ErrCategoryNotFound.Error())
			}

			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postUpdateItem(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		item, ok := readJSON[orderdomain.Item](c)
		if !ok {
			return nil
		}

		if err := orderService.UpdateItem(item); err != nil {
			if errors.Is(err, repo.ErrItemNameTaken) {
				return c.String(http.StatusBadRequest, repo.ErrItemNameTaken.Error())
			}

			if errors.Is(err, repo.ErrItemNameEmpty) {
				return c.String(http.StatusBadRequest, repo.ErrItemNameEmpty.Error())
			}

			if errors.Is(err, repo.ErrItemNotFound) {
				return c.String(http.StatusBadRequest, repo.ErrItemNotFound.Error())
			}

			if errors.Is(err, repo.ErrCategoryNotFound) {
				return c.String(http.StatusBadRequest, repo.ErrCategoryNotFound.Error())
			}

			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postDeleteItem(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param(`id`))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		if err := orderService.DeleteItem(id); err != nil {
			if errors.Is(err, repo.ErrItemNotFound) {
				return c.String(http.StatusBadRequest, repo.ErrItemNotFound.Error())
			}

			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postNewMember(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		member, ok := readJSON[orderdomain.Member](c)
		if !ok {
			return nil
		}

		if member.Club != userFromContext(c).Club {
			return c.String(http.StatusBadRequest, `you cannot only create members for your own club`)
		}

		if err := orderService.NewMember(member); err != nil {
			return allowErrors(c, err,
				repo.ErrMemberNameTaken,
				repo.ErrMemberFieldsNotFilled,
			)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postEditMember(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		member, ok := readJSON[orderdomain.Member](c)
		if !ok {
			return nil
		}

		if err := orderService.EditMember(member); err != nil {
			return allowErrors(c, err,
				repo.ErrMemberNameTaken,
				repo.ErrMemberFieldsNotFilled,
				repo.ErrClubChange,
				repo.ErrMemberNotFound)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postDeleteMember(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param(`id`))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		if err := orderService.DeleteMember(id); err != nil {
			return allowErrors(c, err,
				repo.ErrMemberNotFound,
				repo.ErrMemberHasOrders)
		}

		return c.NoContent(http.StatusOK)
	}
}

func getBillingOrders(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		m, err := orderdomain.ParseMonth(c.QueryParam(`month`))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		orders := orderService.GetOrdersByClub(userFromContext(c).Club, m)

		return c.JSON(http.StatusOK, orders)
	}
}

func getDownload(orderService order.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		m, err := orderdomain.ParseMonth(c.QueryParam(`month`))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		user := userFromContext(c)
		csv := orderService.BillingCSV(user.Club, m)

		filename := m.String() + `-` + user.Club.String() + `.csv`
		c.Response().Header().Set(`content-disposition`, `attachment; filename="`+filename+`"`)
		return c.Blob(http.StatusOK, `text/csv`, csv)
	}
}
