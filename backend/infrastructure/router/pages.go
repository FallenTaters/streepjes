package router

import (
	"net/http"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/templates"
	"go.uber.org/zap"
)

var pageLogger *zap.Logger

type pageData struct {
	ActivePage  string
	User        authdomain.User
	UserDisplay string
	IsBartender bool
	IsAdmin     bool
}

func newPageData(r *http.Request, activePage string) pageData {
	user := userFromContext(r)
	display := user.Username
	if len(display) > 10 {
		display = display[:8] + "…"
	}
	return pageData{
		ActivePage:  activePage,
		User:        user,
		UserDisplay: display,
		IsBartender: user.Role.Has(authdomain.PermissionBarStuff),
		IsAdmin:     user.Role.Has(authdomain.PermissionAdminStuff),
	}
}

func render(w http.ResponseWriter, tmpl string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.Render(w, tmpl, data); err != nil {
		if pageLogger != nil {
			pageLogger.Error("template render failed", zap.String("template", tmpl), zap.Error(err))
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// Profile

func getProfilePage(_ auth.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "profile.html", newPageData(r, "profile"))
	}
}

func postProfilePasswordPage(_ auth.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}

func postProfileNamePage(_ auth.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}

// Bartender pages

func getOrderPage(_ order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "order.html", newPageData(r, "order"))
	}
}

func getHistoryPage(_ order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "history.html", newPageData(r, "history"))
	}
}

func postDeleteOrderPage(_ order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/history", http.StatusSeeOther)
	}
}

func getLeaderboardPage(_ order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "leaderboard.html", newPageData(r, "leaderboard"))
	}
}

// Admin pages

func getUsersPage(_ auth.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "admin/users.html", newPageData(r, "users"))
	}
}

func postUsersPage(_ auth.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	}
}

func postDeleteUserPage(_ auth.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	}
}

func getMembersPage(_ order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "admin/members.html", newPageData(r, "members"))
	}
}

func postMembersPage(_ order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
	}
}

func postDeleteMemberPage(_ order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
	}
}

func getCatalogPage(_ order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "admin/catalog.html", newPageData(r, "catalog"))
	}
}

func postCatalogCategoryPage(_ order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/catalog", http.StatusSeeOther)
	}
}

func postDeleteCategoryPage(_ order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/catalog", http.StatusSeeOther)
	}
}

func postCatalogItemPage(_ order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/catalog", http.StatusSeeOther)
	}
}

func postDeleteItemPage(_ order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/catalog", http.StatusSeeOther)
	}
}

func getBillingPage(_ order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "admin/billing.html", newPageData(r, "billing"))
	}
}
