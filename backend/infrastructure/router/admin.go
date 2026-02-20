package router

import (
	"net/http"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func (s *Server) adminRoutes(mux *http.ServeMux, admin middleware) {
	handle(mux, "GET /admin/users", admin, s.getUsersPage)
	handle(mux, "POST /admin/users", admin, s.postUsersPage)
	handle(mux, "POST /admin/users/{id}/delete", admin, s.postDeleteUserPage)

	handle(mux, "GET /admin/members", admin, s.getMembersPage)
	handle(mux, "POST /admin/members", admin, s.postMembersPage)
	handle(mux, "POST /admin/members/{id}/delete", admin, s.postDeleteMemberPage)

	handle(mux, "GET /admin/catalog", admin, s.getCatalogPage)
	handle(mux, "POST /admin/catalog/category", admin, s.postCatalogCategoryPage)
	handle(mux, "POST /admin/catalog/category/{id}/delete", admin, s.postDeleteCategoryPage)
	handle(mux, "POST /admin/catalog/item", admin, s.postCatalogItemPage)
	handle(mux, "POST /admin/catalog/item/{id}/delete", admin, s.postDeleteItemPage)

	handle(mux, "GET /admin/billing", admin, s.getBillingPage)
	handle(mux, "GET /admin/download", admin, s.getDownload)
}

func (s *Server) getDownload(w http.ResponseWriter, r *http.Request) {
	m, err := orderdomain.ParseMonth(r.URL.Query().Get(`month`))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user := userFromContext(r)
	csv, err := s.order.BillingCSV(user.Club, m)
	if err != nil {
		s.internalError(w, "generate billing CSV", err)
		return
	}

	filename := m.String() + `-` + user.Club.String() + `.csv`
	w.Header().Set(`Content-Disposition`, `attachment; filename="`+filename+`"`)
	w.Header().Set(`Content-Type`, `text/csv`)
	_, _ = w.Write(csv)
}
