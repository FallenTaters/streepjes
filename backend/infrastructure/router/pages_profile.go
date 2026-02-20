package router

import (
	"net/http"

	"github.com/FallenTaters/streepjes/api"
	"go.uber.org/zap"
)

type profileData struct {
	pageData
	PasswordMsg string
	NameMsg     string
}

func (s *Server) getProfilePage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "profile.html", profileData{
		pageData:    newPageData(r, "profile"),
		PasswordMsg: r.URL.Query().Get("pw"),
		NameMsg:     r.URL.Query().Get("name"),
	})
}

func (s *Server) postProfilePasswordPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/profile?pw=error", http.StatusSeeOther)
		return
	}

	user := userFromContext(r)
	err := s.auth.ChangePassword(user, api.ChangePassword{
		Original: r.FormValue("original"),
		New:      r.FormValue("new"),
	})
	if err != nil {
		s.logger.Warn("password change failed", zap.String("user", user.Username), zap.Error(err))
		http.Redirect(w, r, "/profile?pw=error", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile?pw=success", http.StatusSeeOther)
}

func (s *Server) postProfileNamePage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/profile?name=error", http.StatusSeeOther)
		return
	}

	user := userFromContext(r)
	name := r.FormValue("name")

	s.logger.Debug("received change name request",
		zap.String("user", user.Username),
		zap.String("name", name),
	)

	if err := s.auth.ChangeName(user, name); err != nil {
		s.logger.Warn("name change failed", zap.String("user", user.Username), zap.Error(err))
		http.Redirect(w, r, "/profile?name=error", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile?name=success", http.StatusSeeOther)
}
