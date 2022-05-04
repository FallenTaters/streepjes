package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/labstack/echo/v4"
)

func adminRoutes(r *echo.Group, authService auth.Service) {
	r.GET(`/users`, getUsers(authService))
	r.POST(`/users/new`, postNewUser(authService))
	r.POST(`/users/edit`, postEditUser(authService))
	r.POST(`/users/:id/delete`, postDeleteUser(authService))
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
			return c.NoContent(http.StatusBadRequest)
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
