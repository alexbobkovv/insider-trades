package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/config"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/controller/httpapi"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/message"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/repository"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/service"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/httpserver"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/kafka"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/postgresql"

	"github.com/gorilla/mux"
)

func Run(cfg *config.Config) {

	l, err := logger.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.Filepath)
	if err != nil {
		log.Println("app: failed to initialize logger")
	}
	l.Info("app: logger initialized")

	l.Info("postgres url: ", cfg.Postgres.URL)
	psql, err := postgresql.New(cfg.Postgres.URL)
	if err != nil {
		l.Fatalf("app: postgresql.New: %v", err)
	}
	defer psql.Pool.Close()

	kfk, err := kafka.New()
	if err != nil {
		l.Fatalf("app: kafka.New: %v", err)
	}

	insiderTradeService := service.New(
		repository.New(psql),
		message.New(kfk),
	)

	router := mux.NewRouter()

	handler := httpapi.NewHandler(insiderTradeService, l)
	handler.Register(router)

	httpServer := httpserver.New(router, cfg.Server.Port)

	// Starting server with graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		l.Infof("app: server is running on %v", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("app: faild to start server: %v", err)
		}
	}()

	// Waiting for shutdown signal
	sysSignal := <-signalChan
	l.Infof("app: got signal %v, shutting down..", sysSignal)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()

	httpServer.SetKeepAlivesEnabled(false)
	if err := httpServer.Shutdown(ctx); err != nil {
		l.Fatalf("app: %v", err)
	}

	l.Info("app: successful shutdown")
}
