package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArturAda/GO/HW_3/internal/config"
	"github.com/ArturAda/GO/HW_3/internal/controller"
	"github.com/ArturAda/GO/HW_3/internal/service"
	fileStorage "github.com/ArturAda/GO/HW_3/internal/storages/file"
	memoryStorage "github.com/ArturAda/GO/HW_3/internal/storages/memory"
)

func main() {
	cfg := config.CreateNewConfigFrom()
	var storage service.BalanceRepository
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
	svc := service.NewBalanceService(storage)
	srv := controller.New(svc)
	httpSrv := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           srv.GetHandler(),
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
