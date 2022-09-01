package router

import (
	"net/http"
	"strconv"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/go-chi/chi/v5"
)

func bartenderRoutes(r chi.Router, orderService order.Service) {
	r.Get(`/catalog`, getCatalog(orderService))
	r.Get(`/members`, getMembers(orderService))
	r.Get(`/member/{id}`, getMember(orderService))
	r.Get(`/orders`, getOrders(orderService))
	r.Post(`/order`, chio.JSON(postOrder(orderService)))
	r.Post(`/order/{id}/delete`, postDeleteOrder(orderService))
	r.Post(`/leaderboard`, chio.JSON(postGetLeaderboard(orderService)))
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

func getOrders(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r)
		orders := orderService.GetOrdersForBartender(user.ID)
		chio.WriteJSON(w, http.StatusOK, orders)
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

func postDeleteOrder(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bartenderID := userFromContext(r).ID
		orderID, err := strconv.Atoi(chi.URLParam(r, `id`))
		if err != nil {
			chio.WriteString(w, http.StatusUnprocessableEntity, `order id must be integer`)
			return
		}

		ok := orderService.BartenderDeleteOrder(bartenderID, orderID)
		if !ok {
			chio.Empty(w, http.StatusNotFound)
			return
		}

		chio.Empty(w, http.StatusOK)
	}
}

func postGetLeaderboard(orderService order.Service) func(http.ResponseWriter, *http.Request, api.LeaderboardFilter) {
	return func(w http.ResponseWriter, r *http.Request, filter api.LeaderboardFilter) {
		leaderboard := orderService.GetLeaderboard(filter)

		chio.WriteJSON(w, http.StatusOK, leaderboard)
	}
}
