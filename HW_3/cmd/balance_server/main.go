package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ArturAda/GO/HW_3/internal/config"
	"github.com/ArturAda/GO/HW_3/internal/exceptions"
	"github.com/ArturAda/GO/HW_3/internal/functions_and_interface"
	fileStorage "github.com/ArturAda/GO/HW_3/internal/storages/file"
	memoryStorage "github.com/ArturAda/GO/HW_3/internal/storages/memory"
)

type server struct {
	src *functions_and_interface.BalanceService
	mux *http.ServeMux
}

type errResp struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (ser *server) routes() {
	ser.mux.HandleFunc("GET /live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	ser.mux.HandleFunc("GET /api/v1/balance/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		if len(parts) < 5 || parts[3] != "balance" {
			writeJSON(w, http.StatusNotFound, errResp{Error: "not found"})
			return
		}
		userID := parts[len(parts)-1]
		balance, err := ser.src.GetBalance(userID)
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
			"balance_fmt": functions_and_interface.FormatKopeck(balance),
			"currency":    "RUB",
		})
	})
	ser.mux.HandleFunc("POST /api/v1/deposit", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID string `json:"user_id"`
			Amount int64  `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid json"})
			return
		}
		newBalance, err := ser.src.Deposit(req.UserID, req.Amount)
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
			"new_balance_fmt": functions_and_interface.FormatKopeck(newBalance),
		})
	})
	ser.mux.HandleFunc("POST /api/v1/transfer", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			From   string `json:"from"`
			To     string `json:"to"`
			Amount int64  `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid json"})
			return
		}
		fromBalance, toBalance, err := ser.src.Transfer(req.From, req.To, req.Amount)
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
				"balance_fmt": functions_and_interface.FormatKopeck(fromBalance),
			},
			"to": map[string]any{
				"user_id":     req.To,
				"balance":     toBalance,
				"balance_fmt": functions_and_interface.FormatKopeck(toBalance),
			},
		})
	})
}

func main() {
	cfg := config.CreateNewConfigFrom()
	var storage functions_and_interface.BalanceRepository
	switch cfg.StorageType {
	case "file":
		r, err := fileStorage.New(cfg.JSONPath)
		if err != nil {
			log.Fatalf("file storage init: %v", err)
		}
		storage = r
	default:
		storage = memoryStorage.New()
	}
	svc := functions_and_interface.NewBalanceService(storage)
	srv := &server{src: svc, mux: http.NewServeMux()}
	srv.routes()
	httpSrv := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           srv.mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("listening on %s (storage=%s)", cfg.HTTPAddress, cfg.StorageType)
		if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)
	fmt.Println("shut down")
}
