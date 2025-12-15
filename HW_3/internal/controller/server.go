package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ArturAda/GO/HW_3/internal/exceptions"
	"github.com/ArturAda/GO/HW_3/internal/service"
)

type Server struct {
	src *service.BalanceService
	mux *http.ServeMux
}

type errResp struct {
	Error string `json:"error"`
}

func New(src *service.BalanceService) *Server {
	s := &Server{
		src: src,
		mux: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) GetHandler() http.Handler {
	return s.mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	s.mux.HandleFunc("GET /api/v1/balance/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("user_id")
		balance, err := s.src.GetBalance(userID)
		if err != nil {
			switch {
			case errors.Is(err, exceptions.ErrNotFound):
				writeJSON(w, http.StatusNotFound, errResp{Error: "user not found"})
			default:
				writeJSON(w, http.StatusInternalServerError, errResp{Error: err.Error()})
			}
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"user_id":     userID,
			"balance":     balance,
			"balance_fmt": service.FormatKopeck(balance),
			"currency":    "RUB",
		})
	})
	s.mux.HandleFunc("POST /api/v1/deposit", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID string `json:"user_id"`
			Amount int64  `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid json"})
			return
		}
		newBalance, err := s.src.Deposit(req.UserID, req.Amount)
		if err != nil {
			switch {
			case errors.Is(err, exceptions.ErrInvalidAmount):
				writeJSON(w, http.StatusBadRequest, errResp{Error: "amount must be positive"})
			default:
				writeJSON(w, http.StatusInternalServerError, errResp{Error: err.Error()})
			}
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"user_id":         req.UserID,
			"new_balance":     newBalance,
			"new_balance_fmt": service.FormatKopeck(newBalance),
		})
	})
	s.mux.HandleFunc("POST /api/v1/transfer", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			From   string `json:"from"`
			To     string `json:"to"`
			Amount int64  `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid json"})
			return
		}
		fromBalance, toBalance, err := s.src.Transfer(req.From, req.To, req.Amount)
		if err != nil {
			switch {
			case errors.Is(err, exceptions.ErrInvalidAmount):
				writeJSON(w, http.StatusBadRequest, errResp{Error: "amount must be positive"})
			case errors.Is(err, exceptions.ErrSelfTransfer):
				writeJSON(w, http.StatusBadRequest, errResp{Error: "cannot transfer to self"})
			case errors.Is(err, exceptions.ErrNotFound):
				writeJSON(w, http.StatusNotFound, errResp{Error: "source user not found"})
			case errors.Is(err, exceptions.ErrInsufficientFunds):
				writeJSON(w, http.StatusConflict, errResp{Error: "insufficient funds"})
			default:
				writeJSON(w, http.StatusInternalServerError, errResp{Error: err.Error()})
			}
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"from": map[string]any{
				"user_id":     req.From,
				"balance":     fromBalance,
				"balance_fmt": service.FormatKopeck(fromBalance),
			},
			"to": map[string]any{
				"user_id":     req.To,
				"balance":     toBalance,
				"balance_fmt": service.FormatKopeck(toBalance),
			},
		})
	})
}
