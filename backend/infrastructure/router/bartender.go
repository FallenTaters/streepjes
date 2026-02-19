package router

import (
	"net/http"
	"strconv"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func bartenderPageRoutes(r chi.Router, orderService order.Service, logger *zap.Logger) {
	// JSON API endpoints (used by order page JS)
	r.Get(`/api/member/{id}`, getMember(orderService))
	r.Post(`/api/order`, chio.JSON(postOrder(orderService)))

	// Page routes
	r.Get(`/order`, getOrderPage(orderService))
	r.Get(`/history`, getHistoryPage(orderService))
	r.Post(`/history/{id}/delete`, postDeleteOrderPage(orderService))
	r.Get(`/leaderboard`, getLeaderboardPage(orderService))
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
