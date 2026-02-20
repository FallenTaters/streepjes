package router

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"go.uber.org/zap"
)

func (s *Server) bartenderRoutes(mux *http.ServeMux, bar middleware) {
	handle(mux, "GET /api/member/{id}", bar, s.getMember)
	handle(mux, "POST /api/order", bar, s.postOrder)
	handle(mux, "GET /order", bar, s.getOrderPage)
	handle(mux, "GET /history", bar, s.getHistoryPage)
	handle(mux, "POST /history/{id}/delete", bar, s.postDeleteOrderPage)
	handle(mux, "GET /leaderboard", bar, s.getLeaderboardPage)
}

func (s *Server) getMember(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	member, err := s.order.GetMemberDetails(id)
	if err != nil {
		if errors.Is(err, repo.ErrMemberNotFound) {
			http.NotFound(w, r)
			return
		}
		s.internalError(w, "get member details", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(member); err != nil {
		s.logger.Error("encode member response", zap.Error(err))
	}
}

func (s *Server) postOrder(w http.ResponseWriter, r *http.Request) {
	var order orderdomain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.order.PlaceOrder(order, userFromContext(r)); err != nil {
		s.logger.Warn("order placement failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
