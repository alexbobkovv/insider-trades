package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/api-gateway-service/config"
	"github.com/alexbobkovv/insider-trades/api-gateway-service/internal/cache"
	"github.com/alexbobkovv/insider-trades/api-gateway-service/internal/controller/httpapi"
	"github.com/alexbobkovv/insider-trades/api-gateway-service/internal/service"
	"github.com/alexbobkovv/insider-trades/pkg/httpserver"
	"github.com/alexbobkovv/insider-trades/pkg/logger"
	"github.com/alexbobkovv/insider-trades/pkg/redisdb"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.New(os.Getenv("CONFIG_PATH"), os.Getenv("CONFIG_NAME"))
	if err != nil {
		log.Fatalf("main: failed to load server config: %v", err)
	}

	l, err := logger.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.Filepath)
	if err != nil {
		log.Printf("main: failed to initialize zap: %v", err)
	}
	l.Info("main: zap initialized")

	redisClient := redisdb.New(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Username, cfg.Redis.Password)
	redisCache := cache.New(redisClient)

	// Connect to gRPC receiver server
	connToReceiver, err := grpc.Dial(cfg.GRPC.ReceiverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Fatalf("main: failed to dial to gRPC receiver server: %v", err)
	}
	defer func() {
		if err := connToReceiver.Close(); err != nil {
			l.Errorf("main: failed to close connection to gRPC receiver: %v", err)
		}
	}()

	tradeClient := api.NewTradeServiceClient(connToReceiver)

	gatewayService := service.New(tradeClient)

	router := mux.NewRouter()

	handler := httpapi.NewHandler(gatewayService, l, cfg, redisCache)
	handler.Register(router)

	httpServer := httpserver.New(router, cfg.HTTPServer.Port)

	// Starting server with graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	// HTTP(rest) server
	go func() {
		l.Infof("main: http server is running on %v", cfg.HTTPServer.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("main: faild to start http server: %v", err)
		}
	}()

	// Waiting for shutdown signal
	sysSignal := <-signalChan
	l.Infof("main: got signal %v, shutting down..", sysSignal)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()

	httpServer.SetKeepAlivesEnabled(false)

	if err := httpServer.Shutdown(ctx); err != nil {
		l.Fatalf("main: %v", err)
	}

	l.Info("main: successful shutdown")

}
