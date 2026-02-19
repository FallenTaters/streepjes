package router

import (
	"net/http"

	"github.com/FallenTaters/streepjes/api"
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

type profileData struct {
	pageData
	PasswordMsg string
	NameMsg     string
}

func getProfilePage(_ auth.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "profile.html", profileData{
			pageData:    newPageData(r, "profile"),
			PasswordMsg: r.URL.Query().Get("pw"),
			NameMsg:     r.URL.Query().Get("name"),
		})
	}
}

func postProfilePasswordPage(authService auth.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/profile?pw=error", http.StatusSeeOther)
			return
		}

		user := userFromContext(r)
		err := authService.ChangePassword(user, api.ChangePassword{
			Original: r.FormValue("original"),
			New:      r.FormValue("new"),
		})

		if err != nil {
			logger.Warn("password change failed", zap.String("user", user.Username), zap.Error(err))
			http.Redirect(w, r, "/profile?pw=error", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/profile?pw=success", http.StatusSeeOther)
	}
}

func postProfileNamePage(authService auth.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/profile?name=error", http.StatusSeeOther)
			return
		}

		user := userFromContext(r)
		name := r.FormValue("name")

		logger.Debug("received change name request",
			zap.String("user", user.Username),
			zap.String("name", name),
		)

		if err := authService.ChangeName(user, name); err != nil {
			logger.Warn("name change failed", zap.String("user", user.Username), zap.Error(err))
			http.Redirect(w, r, "/profile?name=error", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/profile?name=success", http.StatusSeeOther)
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
