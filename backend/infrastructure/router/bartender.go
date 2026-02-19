package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"go.uber.org/zap"
)

func bartenderPageRoutes(mux *http.ServeMux, bar middleware, orderService order.Service, logger *zap.Logger) {
	handle(mux, "GET /api/member/{id}", bar, getMember(orderService))
	handle(mux, "POST /api/order", bar, postOrder(orderService, logger))

	handle(mux, "GET /order", bar, getOrderPage(orderService, logger))
	handle(mux, "GET /history", bar, getHistoryPage(orderService, logger))
	handle(mux, "POST /history/{id}/delete", bar, postDeleteOrderPage(orderService, logger))
	handle(mux, "GET /leaderboard", bar, getLeaderboardPage(orderService, logger))
}

func getMember(orderService order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		member, ok := orderService.GetMemberDetails(id)
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(member)
	}
}

func postOrder(orderService order.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order orderdomain.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := orderService.PlaceOrder(order, userFromContext(r)); err != nil {
			logger.Warn("order placement failed", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
