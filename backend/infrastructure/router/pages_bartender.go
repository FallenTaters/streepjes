package router

import (
	"encoding/json"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/shared"
	"go.uber.org/zap"
)

type orderData struct {
	pageData
	UserClub    string
	CatalogJSON template.JS
	MembersJSON template.JS
}

func (s *Server) getOrderPage(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r)

	catalog, err := s.order.GetCatalog()
	if err != nil {
		s.internalError(w, "get catalog", err)
		return
	}

	allMembers, err := s.order.GetAllMembers()
	if err != nil {
		s.internalError(w, "get members", err)
		return
	}

	catalogBytes, _ := json.Marshal(catalog)
	membersBytes, _ := json.Marshal(allMembers)

	data := orderData{
		pageData:    newPageData(r, "order"),
		UserClub:    user.Club.String(),
		CatalogJSON: template.JS(catalogBytes),
		MembersJSON: template.JS(membersBytes),
	}

	s.render(w, "order.html", data)
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

func (s *Server) getHistoryPage(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r)

	orders, err := s.order.GetOrdersForBartender(user.ID)
	if err != nil {
		s.internalError(w, "get orders for bartender", err)
		return
	}

	members, err := s.order.GetAllMembers()
	if err != nil {
		s.internalError(w, "get members", err)
		return
	}

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

	s.render(w, "history.html", data)
}

func (s *Server) postDeleteOrderPage(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r)
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := s.order.BartenderDeleteOrder(user.ID, id); err != nil {
		s.logger.Warn("order delete failed", zap.Int("id", id), zap.Error(err))
		http.Redirect(w, r, "/history?error=Unable+to+delete+order.", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/history", http.StatusSeeOther)
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

func (s *Server) getLeaderboardPage(w http.ResponseWriter, r *http.Request) {
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

	leaderboard, err := s.order.GetLeaderboard(api.LeaderboardFilter{
		Start: time.Now().AddDate(-10, 0, 0),
		End:   time.Now().AddDate(10, 0, 0),
	})
	if err != nil {
		s.internalError(w, "get leaderboard", err)
		return
	}

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

	s.render(w, "leaderboard.html", data)
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
