package router

import (
	"net/http"
	"strconv"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func bartenderPageRoutes(r chi.Router, orderService order.Service, logger *zap.Logger) {
	// JSON API endpoints (kept for order page JS)
	r.Get(`/api/catalog`, getCatalog(orderService))
	r.Get(`/api/members`, getMembers(orderService))
	r.Get(`/api/member/{id}`, getMember(orderService))
	r.Post(`/api/order`, chio.JSON(postOrder(orderService)))

	// Page routes
	r.Get(`/order`, getOrderPage(orderService))
	r.Get(`/history`, getHistoryPage(orderService))
	r.Post(`/history/{id}/delete`, postDeleteOrderPage(orderService))
	r.Get(`/leaderboard`, getLeaderboardPage(orderService))
}

func getCatalog(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		catalog := orderService.GetCatalog()
		chio.WriteJSON(w, http.StatusOK, catalog)
	}
}

func getMembers(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		members := orderService.GetAllMembers()
		user := userFromContext(r)

		if user.Role == authdomain.RoleAdmin {
			out := make([]orderdomain.Member, 0, len(members))
			for _, m := range members {
				if m.Club == user.Club {
					out = append(out, m)
				}
			}
			members = out
		}

		chio.WriteJSON(w, http.StatusOK, members)
	}
}

func getMember(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, `id`))
		if err != nil {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		member, ok := orderService.GetMemberDetails(id)
		if !ok {
			chio.Empty(w, http.StatusNotFound)
			return
		}

		chio.WriteJSON(w, http.StatusOK, member)
	}
}

func postOrder(orderService order.Service) func(http.ResponseWriter, *http.Request, orderdomain.Order) {
	return func(w http.ResponseWriter, r *http.Request, order orderdomain.Order) {
		if err := orderService.PlaceOrder(order, userFromContext(r)); err != nil {
			chio.WriteString(w, http.StatusBadRequest, err.Error())
			return
		}

		chio.Empty(w, http.StatusOK)
	}
}
