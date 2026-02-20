package router

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"go.uber.org/zap"
)

type membersData struct {
	pageData

	Members    []orderdomain.Member
	ShowForm   bool
	FormTitle  string
	EditMember *orderdomain.Member
	FormName   string
	Error      string
}

func (s *Server) getMembersPage(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r)

	allMembers, err := s.order.GetAllMembers()
	if err != nil {
		s.internalError(w, "get members", err)
		return
	}

	members := make([]orderdomain.Member, 0, len(allMembers))
	for _, m := range allMembers {
		if m.Club == user.Club {
			members = append(members, m)
		}
	}

	sort.Slice(members, func(i, j int) bool {
		return strings.ToLower(members[i].Name) < strings.ToLower(members[j].Name)
	})

	data := membersData{
		pageData: newPageData(r, "members"),
		Members:  members,
	}

	action := r.URL.Query().Get("action")
	errMsg := r.URL.Query().Get("error")

	switch action {
	case "new": //nolint:goconst // "new" is self-describing as a literal
		data.ShowForm = true
		data.FormTitle = "New Member"
	case "edit":
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		for i := range members {
			if members[i].ID == id {
				m := members[i]
				data.ShowForm = true
				data.EditMember = &m
				data.FormTitle = "Edit " + m.Name
				data.FormName = m.Name
				break
			}
		}
	}

	if errMsg != "" {
		data.Error = errMsg
	}

	s.render(w, "admin/members.html", data)
}

func (s *Server) postMembersPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/admin/members?error=Invalid+form+data", http.StatusSeeOther)
		return
	}

	user := userFromContext(r)
	idStr := r.FormValue("id")
	name := r.FormValue("name")

	if name == "" {
		action := "new"
		if idStr != "" {
			action = "edit&id=" + idStr
		}
		http.Redirect(w, r, "/admin/members?action="+action+"&error=Name+is+required", http.StatusSeeOther)
		return
	}

	if idStr == "" {
		member := orderdomain.Member{
			Club: user.Club,
			Name: name,
		}
		if err := s.order.NewMember(member); err != nil {
			s.logger.Warn("member create failed", zap.Error(err))
			http.Redirect(w, r, "/admin/members?action=new&error=Unable+to+create+member.", http.StatusSeeOther)
			return
		}
	} else {
		id, _ := strconv.Atoi(idStr)
		member := orderdomain.Member{
			ID:   id,
			Club: user.Club,
			Name: name,
		}
		if err := s.order.EditMember(member); err != nil {
			s.logger.Warn("member update failed", zap.Int("id", id), zap.Error(err))
			http.Redirect(w, r, fmt.Sprintf("/admin/members?action=edit&id=%d&error=Unable+to+update+member.", id), http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
}

func (s *Server) postDeleteMemberPage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := s.order.DeleteMember(id); err != nil {
		s.logger.Warn("member delete failed", zap.Int("id", id), zap.Error(err))
		http.Redirect(w, r, fmt.Sprintf("/admin/members?action=edit&id=%d&error=Unable+to+delete+member.", id), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
}
