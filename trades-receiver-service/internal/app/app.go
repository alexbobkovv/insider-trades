package app

import (
	"context"
	"github.com/gorilla/mux"
	"insidertradesreceiver/config"
	"insidertradesreceiver/internal/controller/httpapi"
	"insidertradesreceiver/internal/message"
	"insidertradesreceiver/internal/repository"
	"insidertradesreceiver/internal/service"
	"insidertradesreceiver/pkg/httpserver"
	"insidertradesreceiver/pkg/kafka"
	"insidertradesreceiver/pkg/logger"
	"insidertradesreceiver/pkg/postgresql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config, l *logger.Logger) {

	psql, err := postgresql.New(cfg.Postgres.URL)
	if err != nil {
		l.Fatalf("internal.app.postgresql.New: %v", err)
	}
	defer psql.Pool.Close()

	kafka, err := kafka.New()

	insiderTradeService := service.New(
		repository.New(psql),
		message.New(kafka),
	)

	router := mux.NewRouter()

	handler, err := httpapi.NewHandler(insiderTradeService, l)
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
