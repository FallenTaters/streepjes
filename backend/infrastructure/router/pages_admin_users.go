package router

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"go.uber.org/zap"
)

type usersData struct {
	pageData

	Users          []authdomain.User
	ShowForm       bool
	FormTitle      string
	EditUser       *authdomain.User
	FormUsername   string
	FormName       string
	FormClub       int
	FormRole       int
	FormClubClass  string
	PasswordLabel  string
	SaveButtonText string
	Error          string
}

func clubClass(c domain.Club) string {
	if c == domain.ClubUnknown {
		return "no-club"
	}
	return c.String()
}

func (s *Server) getUsersPage(w http.ResponseWriter, r *http.Request) {
	users, err := s.auth.GetUsers()
	if err != nil {
		s.internalError(w, "get users", err)
		return
	}

	sort.Slice(users, func(i, j int) bool {
		return strings.ToLower(users[i].Name) < strings.ToLower(users[j].Name)
	})

	data := usersData{
		pageData: newPageData(r, "users"),
		Users:    users,
		Error:    r.URL.Query().Get("error"),
	}

	switch r.URL.Query().Get("action") {
	case "new":
		data.ShowForm = true
		data.FormTitle = "New User"
		data.PasswordLabel = "Password"
		data.SaveButtonText = "Add User"
		data.FormClubClass = "no-club"
	case "edit":
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)
		for i := range users {
			if users[i].ID == id {
				u := users[i]
				data.ShowForm = true
				data.EditUser = &u
				data.FormTitle = "Edit User: " + u.Name
				data.FormUsername = u.Username
				data.FormName = u.Name
				data.FormClub = int(u.Club)
				data.FormRole = int(u.Role)
				data.FormClubClass = clubClass(u.Club)
				data.PasswordLabel = "New Password (optional)"
				data.SaveButtonText = "Save Changes"
				break
			}
		}
	}

	s.render(w, "admin/users.html", data)
}

func (s *Server) postUsersPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/admin/users?error=Invalid+form+data", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	username := r.FormValue("username")
	password := r.FormValue("password")
	name := r.FormValue("name")
	clubInt, _ := strconv.Atoi(r.FormValue("club"))
	roleInt, _ := strconv.Atoi(r.FormValue("role"))

	club := domain.Club(clubInt)
	role := authdomain.Role(roleInt)

	if idStr == "" {
		if username == "" || password == "" || name == "" || club == domain.ClubUnknown || role == authdomain.RoleNotAuthorized {
			http.Redirect(w, r, "/admin/users?action=new&error=All+fields+must+be+filled", http.StatusSeeOther)
			return
		}

		user := authdomain.User{
			Username: username,
			Name:     name,
			Club:     club,
			Role:     role,
		}

		if err := s.auth.Register(user, password); err != nil {
			s.logger.Warn("user create failed", zap.Error(err))
			http.Redirect(w, r, "/admin/users?action=new&error=Unable+to+create+user.+Maybe+the+username+is+already+taken.", http.StatusSeeOther)
			return
		}
	} else {
		id, _ := strconv.Atoi(idStr)

		if username == "" || name == "" || club == domain.ClubUnknown || role == authdomain.RoleNotAuthorized {
			http.Redirect(w, r, fmt.Sprintf("/admin/users?action=edit&id=%d&error=All+fields+must+be+filled", id), http.StatusSeeOther)
			return
		}

		user := authdomain.User{
			ID:       id,
			Username: username,
			Name:     name,
			Club:     club,
			Role:     role,
		}

		if err := s.auth.Update(user, password); err != nil {
			s.logger.Warn("user update failed", zap.Int("id", id), zap.Error(err))
			http.Redirect(w, r, fmt.Sprintf("/admin/users?action=edit&id=%d&error=Unable+to+update+user.+Maybe+the+username+is+already+taken.", id), http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (s *Server) postDeleteUserPage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := s.auth.Delete(id); err != nil {
		s.logger.Warn("user delete failed", zap.Int("id", id), zap.Error(err))
		http.Redirect(w, r, fmt.Sprintf("/admin/users?action=edit&id=%d&error=Unable+to+delete+user.", id), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
