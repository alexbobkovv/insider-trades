package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/rabbitmq"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/config"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/controller/grpcapi"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/controller/httpapi"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/message"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/repository"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/service"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/httpserver"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/postgresql"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config) {

	l, err := logger.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.Filepath)
	if err != nil {
		log.Printf("app: failed to initialize zap: %v", err)
	}
	l.Info("app: zap initialized")

	psql, err := postgresql.New(cfg.Postgres.URL)
	if err != nil {
		l.Fatalf("app: postgresql.New: %v", err)
	}
	defer psql.Pool.Close()

	rmq, err := rabbitmq.New(cfg.AmqpURL)
	defer func() {
		err := rmq.Channel.Close()
		if err != nil {
			l.Errorf("app: failed to close amqp server channel")
		}
	}()
	if err != nil {
		l.Fatalf("app: failed to connect to RabbitMQ: %v", err)
	}

	messageBroker, err := message.New(rmq, cfg.RabbitMQ)
	if err != nil {
		l.Fatalf("app: failed to initialize messageBroker: %v", err)
	}

	insiderTradeService := service.New(
		repository.New(psql),
		messageBroker,
	)

	router := mux.NewRouter()

	handler := httpapi.NewHandler(insiderTradeService, l, cfg)
	handler.Register(router)

	httpServer := httpserver.New(router, cfg.HTTPServer.Port)

	// Setup gRPC server
	lis, err := net.Listen("tcp", cfg.GRPCServer.Port)
	if err != nil {
		l.Fatalf("app: failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcTradeServer := grpcapi.NewTradeServer(insiderTradeService, l, cfg)
	api.RegisterTradeServiceServer(grpcServer, grpcTradeServer)

	// Starting servers with graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	// HTTP(rest) server
	go func() {
		l.Infof("app: http server is running on %v", cfg.HTTPServer.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("app: faild to start http server: %v", err)
		}
	}()

	// gRPC server
	go func() {
		l.Infof("app: gRPC server is running on %v", cfg.GRPCServer.Port)
		if err := grpcServer.Serve(lis); err != nil {
			l.Fatalf("app: faild to start gRPC server: %v", err)
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
	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		l.Fatalf("app: %v", err)
	}

	l.Info("app: successful shutdown")
}
