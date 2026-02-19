package router

import (
	"net/http"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"go.uber.org/zap"
)

func adminPageRoutes(mux *http.ServeMux, admin middleware, authService auth.Service, orderService order.Service, logger *zap.Logger) {
	handle(mux, "GET /admin/users", admin, getUsersPage(authService, logger))
	handle(mux, "POST /admin/users", admin, postUsersPage(authService, logger))
	handle(mux, "POST /admin/users/{id}/delete", admin, postDeleteUserPage(authService, logger))

	handle(mux, "GET /admin/members", admin, getMembersPage(orderService, logger))
	handle(mux, "POST /admin/members", admin, postMembersPage(orderService, logger))
	handle(mux, "POST /admin/members/{id}/delete", admin, postDeleteMemberPage(orderService, logger))

	handle(mux, "GET /admin/catalog", admin, getCatalogPage(orderService, logger))
	handle(mux, "POST /admin/catalog/category", admin, postCatalogCategoryPage(orderService, logger))
	handle(mux, "POST /admin/catalog/category/{id}/delete", admin, postDeleteCategoryPage(orderService, logger))
	handle(mux, "POST /admin/catalog/item", admin, postCatalogItemPage(orderService, logger))
	handle(mux, "POST /admin/catalog/item/{id}/delete", admin, postDeleteItemPage(orderService, logger))

	handle(mux, "GET /admin/billing", admin, getBillingPage(orderService, logger))
	handle(mux, "GET /admin/download", admin, getDownload(orderService, logger))
}

func getDownload(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := orderdomain.ParseMonth(r.URL.Query().Get(`month`))
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user := userFromContext(r)
		csv := orderService.BillingCSV(user.Club, m)

		filename := m.String() + `-` + user.Club.String() + `.csv`
		w.Header().Set(`Content-Disposition`, `attachment; filename="`+filename+`"`)
		w.Header().Set(`Content-Type`, `text/csv`)
		w.Write(csv)
	}
}
