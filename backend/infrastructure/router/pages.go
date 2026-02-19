package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/shared"
	"github.com/FallenTaters/streepjes/templates"
	"github.com/go-chi/chi/v5"
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

type orderData struct {
	pageData
	UserClub    string
	CatalogJSON template.JS
	MembersJSON template.JS
}

func getOrderPage(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r)
		catalog := orderService.GetCatalog()
		allMembers := orderService.GetAllMembers()

		catalogBytes, _ := json.Marshal(catalog)
		membersBytes, _ := json.Marshal(allMembers)

		data := orderData{
			pageData:    newPageData(r, "order"),
			UserClub:    user.Club.String(),
			CatalogJSON: template.JS(catalogBytes),
			MembersJSON: template.JS(membersBytes),
		}

		render(w, "order.html", data)
	}
}

type historyLine struct {
	Amount    int
	Name      string
	LinePrice orderdomain.Price
}

type historyOrder struct {
	ID         int
	Date       string
	MemberName string
	ClubName   string
	ClubClass  string
	Price      orderdomain.Price
	Lines      []historyLine
	Cancelled  bool
}

type historyData struct {
	pageData
	Orders []historyOrder
	Error  string
}

func getHistoryPage(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r)
		orders := orderService.GetOrdersForBartender(user.ID)
		members := orderService.GetAllMembers()

		membersByID := make(map[int]orderdomain.Member)
		for _, m := range members {
			membersByID[m.ID] = m
		}

		sort.Slice(orders, func(i, j int) bool {
			return orders[i].OrderTime.After(orders[j].OrderTime)
		})

		viewOrders := make([]historyOrder, 0, len(orders))
		for _, o := range orders {
			member := membersByID[o.MemberID]

			var lines []historyLine
			var parsed []orderdomain.Line
			if err := json.Unmarshal([]byte(o.Contents), &parsed); err == nil {
				for _, l := range parsed {
					lines = append(lines, historyLine{
						Amount:    l.Amount,
						Name:      l.Item.Name,
						LinePrice: l.Price(o.Club),
					})
				}
			}

			cc := o.Club.String()
			if o.Status == orderdomain.StatusCancelled {
				cc = "grey"
			}

			viewOrders = append(viewOrders, historyOrder{
				ID:         o.ID,
				Date:       shared.PrettyDatetime(o.OrderTime),
				MemberName: member.Name,
				ClubName:   member.Club.String(),
				ClubClass:  cc,
				Price:      o.Price,
				Lines:      lines,
				Cancelled:  o.Status == orderdomain.StatusCancelled,
			})
		}

		data := historyData{
			pageData: newPageData(r, "history"),
			Orders:   viewOrders,
			Error:    r.URL.Query().Get("error"),
		}

		render(w, "history.html", data)
	}
}

func postDeleteOrderPage(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r)
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		if err := orderService.BartenderDeleteOrder(user.ID, id); err != nil {
			pageLogger.Warn("order delete failed", zap.Int("id", id), zap.Error(err))
			http.Redirect(w, r, "/history?error=Unable+to+delete+order.", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/history", http.StatusSeeOther)
	}
}

var itemWeights = map[string]int{
	"Bier":           1,
	"Weizen glas":    1,
	"Pitcher":        5,
	"Weizen Pitcher": 5,
	"Flugel":         1,
	"Bier Barcie":    1,
	"Bier BarCie":    1,
	"Seltzer BarCie": 1,
	"Seltzer":        1,
	"Wine Bottle":    5,
	"Wijn":           1,
	"Radler":         0,
}

type leaderboardRank struct {
	Name      string
	ClubClass string
	Total     string
	Items     []string
}

type leaderboardData struct {
	pageData
	Gladiators bool
	Parabool   bool
	Calamari   bool
	Sort       string
	Total      string
	Ranking    []leaderboardRank
}

func getLeaderboardPage(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		gladiators := q.Get("gladiators") == "1"
		parabool := q.Get("parabool") == "1"
		calamari := q.Get("calamari") == "1"
		sortMode := q.Get("sort")

		hasFilters := q.Has("gladiators") || q.Has("parabool") || q.Has("calamari")
		if !hasFilters {
			gladiators = true
			parabool = true
			calamari = true
		}
		if sortMode == "" {
			sortMode = "items"
		}

		leaderboard := orderService.GetLeaderboard(api.LeaderboardFilter{
			Start: time.Now().AddDate(-10, 0, 0),
			End:   time.Now().AddDate(10, 0, 0),
		})

		var totalStr string
		var ranking []api.LeaderboardRank

		if sortMode == "money" {
			total, r := leaderboard.MoneyRanking()
			totalStr = total.String()
			ranking = r
		} else {
			total, r := leaderboard.ItemRanking(itemWeights)
			totalStr = strconv.Itoa(total)
			ranking = r
		}

		filtered := make([]leaderboardRank, 0, len(ranking))
		for _, rank := range ranking {
			show := (rank.Club == domain.ClubGladiators && gladiators) ||
				(rank.Club == domain.ClubParabool && parabool) ||
				(rank.Club == domain.ClubCalamari && calamari)
			if !show {
				continue
			}

			items := sortItemInfo(rank.ItemInfo)
			filtered = append(filtered, leaderboardRank{
				Name:      rank.Name,
				ClubClass: rank.Club.String(),
				Total:     rank.Total,
				Items:     items,
			})
		}

		data := leaderboardData{
			pageData:   newPageData(r, "leaderboard"),
			Gladiators: gladiators,
			Parabool:   parabool,
			Calamari:   calamari,
			Sort:       sortMode,
			Total:      totalStr,
			Ranking:    filtered,
		}

		render(w, "leaderboard.html", data)
	}
}

func sortItemInfo(itemInfo map[string]int) []string {
	type mc struct {
		msg   string
		count int
	}

	out := make([]mc, 0, len(itemInfo))
	for name, count := range itemInfo {
		if w, ok := itemWeights[name]; !ok || w > 0 {
			out = append(out, mc{
				msg:   strconv.Itoa(count) + " " + name,
				count: count,
			})
		}
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].count == out[j].count {
			return out[i].msg < out[j].msg
		}
		return out[i].count > out[j].count
	})

	result := make([]string, len(out))
	for i, o := range out {
		result[i] = o.msg
	}
	return result
}

// Admin pages

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

func getUsersPage(authService auth.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := authService.GetUsers()
		sort.Slice(users, func(i, j int) bool {
			return strings.ToLower(users[i].Name) < strings.ToLower(users[j].Name)
		})

		data := usersData{
			pageData: newPageData(r, "users"),
			Users:    users,
		}

		action := r.URL.Query().Get("action")
		errMsg := r.URL.Query().Get("error")

		switch action {
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

		if errMsg != "" {
			data.Error = errMsg
		}

		render(w, "admin/users.html", data)
	}
}

func postUsersPage(authService auth.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

			if err := authService.Register(user, password); err != nil {
				logger.Warn("user create failed", zap.Error(err))
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

			if err := authService.Update(user, password); err != nil {
				logger.Warn("user update failed", zap.Int("id", id), zap.Error(err))
				http.Redirect(w, r, fmt.Sprintf("/admin/users?action=edit&id=%d&error=Unable+to+update+user.+Maybe+the+username+is+already+taken.", id), http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	}
}

func postDeleteUserPage(authService auth.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		if err := authService.Delete(id); err != nil {
			logger.Warn("user delete failed", zap.Int("id", id), zap.Error(err))
			http.Redirect(w, r, fmt.Sprintf("/admin/users?action=edit&id=%d&error=Unable+to+delete+user.", id), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	}
}

type membersData struct {
	pageData
	Members    []orderdomain.Member
	ShowForm   bool
	FormTitle  string
	EditMember *orderdomain.Member
	FormName   string
	Error      string
}

func getMembersPage(orderService order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r)
		allMembers := orderService.GetAllMembers()

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
		case "new":
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

		render(w, "admin/members.html", data)
	}
}

func postMembersPage(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			if err := orderService.NewMember(member); err != nil {
				logger.Warn("member create failed", zap.Error(err))
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
			if err := orderService.EditMember(member); err != nil {
				logger.Warn("member update failed", zap.Int("id", id), zap.Error(err))
				http.Redirect(w, r, fmt.Sprintf("/admin/members?action=edit&id=%d&error=Unable+to+update+member.", id), http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
	}
}

func postDeleteMemberPage(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		if err := orderService.DeleteMember(id); err != nil {
			logger.Warn("member delete failed", zap.Int("id", id), zap.Error(err))
			http.Redirect(w, r, fmt.Sprintf("/admin/members?action=edit&id=%d&error=Unable+to+delete+member.", id), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/admin/members", http.StatusSeeOther)
	}
}

type catalogData struct {
	pageData
	Categories         []orderdomain.Category
	DisplayItems       []orderdomain.Item
	SelectedCategory   *orderdomain.Category
	SelectedCategoryID int
	SelectedItemID     int
	FormTitle          string

	ShowCategoryForm bool
	EditCategory     *orderdomain.Category
	CategoryName     string

	ShowItemForm    bool
	EditItem        *orderdomain.Item
	ItemCategoryID  int
	ItemName        string
	PriceGladiators int
	PriceParabool   int
	PriceCalamari   int

	Error string
}

func getCatalogPage(orderService order.Service, _ *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		catalog := orderService.GetCatalog()

		categories := catalog.Categories
		items := catalog.Items

		sort.Slice(categories, func(i, j int) bool {
			return strings.ToLower(categories[i].Name) < strings.ToLower(categories[j].Name)
		})
		sort.Slice(items, func(i, j int) bool {
			return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		})

		data := catalogData{
			pageData:   newPageData(r, "catalog"),
			Categories: categories,
		}

		catID, _ := strconv.Atoi(r.URL.Query().Get("cat"))
		if catID != 0 {
			for i := range categories {
				if categories[i].ID == catID {
					c := categories[i]
					data.SelectedCategory = &c
					data.SelectedCategoryID = c.ID
					data.FormTitle = "Edit Category - " + c.Name
					data.ShowCategoryForm = true
					data.EditCategory = &c
					data.CategoryName = c.Name
					break
				}
			}

			displayItems := make([]orderdomain.Item, 0)
			for _, item := range items {
				if item.CategoryID == catID {
					displayItems = append(displayItems, item)
				}
			}
			data.DisplayItems = displayItems
		}

		action := r.URL.Query().Get("action")
		errMsg := r.URL.Query().Get("error")

		switch action {
		case "new-category":
			data.ShowCategoryForm = true
			data.ShowItemForm = false
			data.EditCategory = nil
			data.FormTitle = "New Category"
			data.CategoryName = ""
		case "edit-item":
			itemID, _ := strconv.Atoi(r.URL.Query().Get("id"))
			for i := range items {
				if items[i].ID == itemID {
					it := items[i]
					data.ShowCategoryForm = false
					data.ShowItemForm = true
					data.EditItem = &it
					data.SelectedItemID = it.ID
					data.FormTitle = "Edit Item - " + it.Name
					data.ItemCategoryID = it.CategoryID
					data.ItemName = it.Name
					data.PriceGladiators = int(it.PriceGladiators)
					data.PriceParabool = int(it.PriceParabool)
					data.PriceCalamari = int(it.PriceCalamari)
					break
				}
			}
		case "new-item":
			data.ShowCategoryForm = false
			data.ShowItemForm = true
			data.FormTitle = "New Item"
			if data.SelectedCategory != nil {
				data.ItemCategoryID = data.SelectedCategory.ID
			}
		}

		if errMsg != "" {
			data.Error = errMsg
		}

		render(w, "admin/catalog.html", data)
	}
}

func catalogRedirect(catID int, errMsg string) string {
	u := "/admin/catalog"
	if catID != 0 {
		u += fmt.Sprintf("?cat=%d", catID)
	}
	if errMsg != "" {
		sep := "?"
		if catID != 0 {
			sep = "&"
		}
		u += sep + "error=" + errMsg
	}
	return u
}

func postCatalogCategoryPage(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, catalogRedirect(0, "Invalid+form+data"), http.StatusSeeOther)
			return
		}

		idStr := r.FormValue("id")
		name := r.FormValue("name")

		if idStr == "" {
			cat := orderdomain.Category{Name: name}
			if err := orderService.NewCategory(cat); err != nil {
				logger.Warn("category create failed", zap.Error(err))
				http.Redirect(w, r, catalogRedirect(0, "Unable+to+create+category."), http.StatusSeeOther)
				return
			}
		} else {
			id, _ := strconv.Atoi(idStr)
			cat := orderdomain.Category{ID: id, Name: name}
			if err := orderService.UpdateCategory(cat); err != nil {
				logger.Warn("category update failed", zap.Int("id", id), zap.Error(err))
				http.Redirect(w, r, catalogRedirect(id, "Unable+to+update+category."), http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, "/admin/catalog", http.StatusSeeOther)
	}
}

func postDeleteCategoryPage(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		if err := orderService.DeleteCategory(id); err != nil {
			logger.Warn("category delete failed", zap.Int("id", id), zap.Error(err))
			http.Redirect(w, r, catalogRedirect(id, "Unable+to+delete+category.+It+may+still+have+items."), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/admin/catalog", http.StatusSeeOther)
	}
}

func postCatalogItemPage(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, catalogRedirect(0, "Invalid+form+data"), http.StatusSeeOther)
			return
		}

		idStr := r.FormValue("id")
		catID, _ := strconv.Atoi(r.FormValue("category_id"))
		name := r.FormValue("name")
		priceGlad, _ := strconv.Atoi(r.FormValue("price_gladiators"))
		pricePara, _ := strconv.Atoi(r.FormValue("price_parabool"))
		priceCala, _ := strconv.Atoi(r.FormValue("price_calamari"))

		if idStr == "" {
			item := orderdomain.Item{
				CategoryID:      catID,
				Name:            name,
				PriceGladiators: orderdomain.Price(priceGlad),
				PriceParabool:   orderdomain.Price(pricePara),
				PriceCalamari:   orderdomain.Price(priceCala),
			}
			if err := orderService.NewItem(item); err != nil {
				logger.Warn("item create failed", zap.Error(err))
				http.Redirect(w, r, catalogRedirect(catID, "Unable+to+create+item."), http.StatusSeeOther)
				return
			}
		} else {
			id, _ := strconv.Atoi(idStr)
			item := orderdomain.Item{
				ID:              id,
				CategoryID:      catID,
				Name:            name,
				PriceGladiators: orderdomain.Price(priceGlad),
				PriceParabool:   orderdomain.Price(pricePara),
				PriceCalamari:   orderdomain.Price(priceCala),
			}
			if err := orderService.UpdateItem(item); err != nil {
				logger.Warn("item update failed", zap.Int("id", id), zap.Error(err))
				http.Redirect(w, r, catalogRedirect(catID, "Unable+to+update+item."), http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, catalogRedirect(catID, ""), http.StatusSeeOther)
	}
}

func postDeleteItemPage(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		catIDStr := r.URL.Query().Get("cat")
		catID, _ := strconv.Atoi(catIDStr)

		if err := orderService.DeleteItem(id); err != nil {
			logger.Warn("item delete failed", zap.Int("id", id), zap.Error(err))
			http.Redirect(w, r, catalogRedirect(catID, "Unable+to+delete+item."), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, catalogRedirect(catID, ""), http.StatusSeeOther)
	}
}

type billingOrder struct {
	MemberName string
	Price      orderdomain.Price
	OrderTime  string
	Lines      []string
}

type billingData struct {
	pageData
	Month  string
	Orders []billingOrder
}

func parseOrderLines(contents string) []string {
	var lines []orderdomain.Line
	if err := json.Unmarshal([]byte(contents), &lines); err != nil {
		return []string{"order data unreadable"}
	}

	out := make([]string, 0, len(lines))
	for _, line := range lines {
		out = append(out, fmt.Sprintf("%dx %s", line.Amount, line.Item.Name))
	}
	return out
}

func getBillingPage(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r)

		data := billingData{
			pageData: newPageData(r, "billing"),
		}

		monthStr := r.URL.Query().Get("month")
		if monthStr != "" {
			month, err := orderdomain.ParseMonth(monthStr)
			if err == nil {
				data.Month = month.String()

				orders := orderService.GetOrdersByClub(user.Club, month)
				members := orderService.GetAllMembers()

				membersByID := make(map[int]orderdomain.Member)
				for _, m := range members {
					membersByID[m.ID] = m
				}

				sort.Slice(orders, func(i, j int) bool {
					name1 := strings.ToLower(membersByID[orders[i].MemberID].Name)
					name2 := strings.ToLower(membersByID[orders[j].MemberID].Name)
					if name1 == name2 {
						return orders[i].OrderTime.Before(orders[j].OrderTime)
					}
					return name1 < name2
				})

				billingOrders := make([]billingOrder, 0, len(orders))
				for _, o := range orders {
					if o.MemberID == 0 {
						continue
					}
					billingOrders = append(billingOrders, billingOrder{
						MemberName: membersByID[o.MemberID].Name,
						Price:      o.Price,
						OrderTime:  o.OrderTime.Format("2006-01-02 15:04"),
						Lines:      parseOrderLines(o.Contents),
					})
				}
				data.Orders = billingOrders
			}
		}

		render(w, "admin/billing.html", data)
	}
}
