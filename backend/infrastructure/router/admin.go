package router

import (
	"net/http"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func adminPageRoutes(r chi.Router, authService auth.Service, orderService order.Service, logger *zap.Logger) {
	r.Get(`/users`, getUsersPage(authService, logger))
	r.Post(`/users`, postUsersPage(authService, logger))
	r.Post(`/users/{id}/delete`, postDeleteUserPage(authService, logger))

	r.Get(`/members`, getMembersPage(orderService, logger))
	r.Post(`/members`, postMembersPage(orderService, logger))
	r.Post(`/members/{id}/delete`, postDeleteMemberPage(orderService, logger))

	r.Get(`/catalog`, getCatalogPage(orderService, logger))
	r.Post(`/catalog/category`, postCatalogCategoryPage(orderService, logger))
	r.Post(`/catalog/category/{id}/delete`, postDeleteCategoryPage(orderService, logger))
	r.Post(`/catalog/item`, postCatalogItemPage(orderService, logger))
	r.Post(`/catalog/item/{id}/delete`, postDeleteItemPage(orderService, logger))

	r.Get(`/billing`, getBillingPage(orderService))
	r.Get(`/download`, getDownload(orderService))
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
