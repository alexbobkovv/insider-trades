package app

import (
	"context"
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

func Run(cfg *config.Config, l *logger.Logger) {

	psql, err := postgresql.New(cfg.Postgres.URL)
	if err != nil {
		l.Fatalf("internal.app.postgresql.New: %v", err)
	}
	defer psql.Pool.Close()

	kafka, err := kafka.New()
	if err != nil {
		l.Fatalf("internal.app.kafka.New: %v", err)
	}

	insiderTradeService := service.New(
		repository.New(psql, l),
		message.New(kafka),
		l,
	)

	router := mux.NewRouter()

	handler, err := httpapi.NewHandler(insiderTradeService, l)
	if err != nil {
		l.Fatalf("internal.app.httpapi.NewHandler: %v", err)
	}

	handler.Register(router)

	httpServer := httpserver.New(router, cfg.Server.Port)

	// Starting server with graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		l.Infof("Server is running on %v", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("Faild to start server: %v", err)
		}
	}()

	// Waiting for shutdown signal
	sysSignal := <-signalChan
	l.Infof("Got signal %v, shutting down..", sysSignal)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer func() {
		cancel()
	}()

	httpServer.SetKeepAlivesEnabled(false)

	if err := httpServer.Shutdown(ctx); err != nil {
		l.Fatalf("Error: %v", err)
	}

	l.Info("Successful shutdown")

}
