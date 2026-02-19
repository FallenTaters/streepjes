package router

// TODO: this package should not depend on repo. move errors to application layer (will require large refactor)
// better idea: filter errors in application layer, this layer should not be worried about it.
// maybe make a type of error that can be returned, that way the router package can simply check if the error is of that type, and if so, write it.

import (
	"net/http"
	"strconv"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func adminRoutes(r chi.Router, authService auth.Service, orderService order.Service, logger *zap.Logger) {
	r.Get(`/users`, getUsers(authService))
	r.Post(`/users/new`, postNewUser(authService, logger))
	r.Post(`/users/edit`, postEditUser(authService, logger))
	r.Post(`/users/{id}/delete`, postDeleteUser(authService))

	r.Post(`/category/new`, postNewCategory(orderService, logger))
	r.Post(`/category/update`, postUpdateCategory(orderService, logger))
	r.Post(`/category/{id}/delete`, postDeleteCategory(orderService, logger))

	r.Post(`/item/new`, postNewItem(orderService, logger))
	r.Post(`/item/update`, postUpdateItem(orderService, logger))
	r.Post(`/item/{id}/delete`, postDeleteItem(orderService, logger))

	r.Post(`/members/new`, postNewMember(orderService, logger))
	r.Post(`/members/edit`, postEditMember(orderService, logger))
	r.Post(`/members/{id}/delete`, postDeleteMember(orderService, logger))

	r.Get(`/billing/orders`, getBillingOrders(orderService))
	r.Get(`/download`, getDownload(orderService))
}

func getUsers(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := authService.GetUsers()
		chio.WriteJSON(w, http.StatusOK, users)
	}
}

func postNewUser(authService auth.Service, logger *zap.Logger) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, user api.UserWithPassword) {
		err := authService.Register(user.User, user.Password)
		if err != nil {
			allowErrors(w, logger, err,
				repo.ErrUsernameTaken,
				repo.ErrUserMissingFields,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}

func postEditUser(authService auth.Service, logger *zap.Logger) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, user api.UserWithPassword) {
		err := authService.Update(user.User, user.Password)
		if err != nil {
			allowErrors(w, logger, err,
				repo.ErrUserNotFound,
				repo.ErrUsernameTaken,
				repo.ErrUserMissingFields,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}

func postDeleteUser(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, `id`))
		if err != nil {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		ok := authService.Delete(id)
		if !ok {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		chio.Empty(w, http.StatusOK)
	}
}

func postNewCategory(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, cat orderdomain.Category) {
		if err := orderService.NewCategory(cat); err != nil {
			allowErrors(w, logger, err,
				repo.ErrCategoryNameTaken,
				repo.ErrCategoryNameEmpty,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}

func postUpdateCategory(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, cat orderdomain.Category) {
		if err := orderService.UpdateCategory(cat); err != nil {
			allowErrors(w, logger, err,
				repo.ErrCategoryNameTaken,
				repo.ErrCategoryNameEmpty,
				repo.ErrCategoryNotFound,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}

func postDeleteCategory(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, `id`))
		if err != nil {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		if err := orderService.DeleteCategory(id); err != nil {
			allowErrors(w, logger, err,
				repo.ErrCategoryNotFound,
				repo.ErrCategoryHasItems,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	}
}

func postNewItem(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, item orderdomain.Item) {
		if err := orderService.NewItem(item); err != nil {
			allowErrors(w, logger, err,
				repo.ErrItemNameTaken,
				repo.ErrItemNameEmpty,
				repo.ErrCategoryNotFound,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}

func postUpdateItem(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, item orderdomain.Item) {
		if err := orderService.UpdateItem(item); err != nil {
			allowErrors(w, logger, err,
				repo.ErrItemNameTaken,
				repo.ErrItemNameEmpty,
				repo.ErrItemNotFound,
				repo.ErrCategoryNotFound,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}

func postDeleteItem(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, `id`))
		if err != nil {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		if err := orderService.DeleteItem(id); err != nil {
			allowErrors(w, logger, err, repo.ErrItemNotFound)
			return
		}

		chio.Empty(w, http.StatusOK)
	}
}

func postNewMember(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, member orderdomain.Member) {
		if member.Club != userFromContext(r).Club {
			chio.WriteString(w, http.StatusBadRequest, `you can only create members for your own club`)
			return
		}

		if err := orderService.NewMember(member); err != nil {
			allowErrors(w, logger, err,
				repo.ErrMemberNameTaken,
				repo.ErrMemberFieldsNotFilled,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}

func postEditMember(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, member orderdomain.Member) {
		if err := orderService.EditMember(member); err != nil {
			allowErrors(w, logger, err,
				repo.ErrMemberNameTaken,
				repo.ErrMemberFieldsNotFilled,
				repo.ErrClubChange,
				repo.ErrMemberNotFound,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}

func postDeleteMember(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, `id`))
		if err != nil {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		if err := orderService.DeleteMember(id); err != nil {
			allowErrors(w, logger, err,
				repo.ErrMemberNotFound,
				repo.ErrMemberHasOrders,
			)
			return
		}

		chio.Empty(w, http.StatusOK)
	}
}

func getBillingOrders(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := orderdomain.ParseMonth(r.URL.Query().Get(`month`))
		if err != nil {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		orders := orderService.GetOrdersByClub(userFromContext(r).Club, m)

		chio.WriteJSON(w, http.StatusOK, orders)
	}
}

func getDownload(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := orderdomain.ParseMonth(r.URL.Query().Get(`month`))
		if err != nil {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		user := userFromContext(r)
		csv := orderService.BillingCSV(user.Club, m)

		filename := m.String() + `-` + user.Club.String() + `.csv`
		w.Header().Set(`content-disposition`, `attachment; filename="`+filename+`"`)
		chio.WriteBlob(w, http.StatusOK, `text/csv`, csv)
	}
}
