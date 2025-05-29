package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"taskManager/statistics/config"
	grpcserver "taskManager/statistics/internal/adapter/grpc/service"
	httpserver "taskManager/statistics/internal/adapter/http/service"
	natshandler "taskManager/statistics/internal/adapter/nats/handler"
	"taskManager/statistics/internal/adapter/postgres"
	"taskManager/statistics/internal/usecase"
	natsconn "taskManager/statistics/pkg/nats"
	natsconsumer "taskManager/statistics/pkg/nats/consumer"
	postgreconn "taskManager/statistics/pkg/postgre"
)

const serviceName = "statistics-service"

type App struct {
	httpServer         *httpserver.API
	grpcServer         *grpcserver.API
	natsPubSubConsumer *natsconsumer.PubSub
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Println(fmt.Sprintf("starting %v service", serviceName))

	log.Println("connecting to postgres", "database", cfg.Postgres.DatabaseName)
	postgresDB, err := postgreconn.NewDB(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	// nats client
	log.Println("connecting to NATS", "hosts", strings.Join(cfg.Nats.Hosts, ","))
	natsClient, err := natsconn.NewClient(ctx, cfg.Nats.Hosts, cfg.Nats.NKey, cfg.Nats.IsTest)
	if err != nil {
		return nil, fmt.Errorf("nats.NewClient: %w", err)
	}
	log.Println("NATS connection status is", natsClient.Conn.Status().String())

	natsPubSubConsumer := natsconsumer.NewPubSub(natsClient)

	postgresDB.Migrate()

	// Repository
	statisticsRepo := postgres.New(postgresDB)

	// Usecase
	statisticsUsecase := usecase.New(statisticsRepo)

	// Nats handler
	natsHandler := natshandler.NewClient(statisticsUsecase)
	natsPubSubConsumer.Subscribe(natsconsumer.PubSubSubscriptionConfig{
		Subject: cfg.Nats.NatsSubjects.UserCreatedEventSubject,
		Handler: natsHandler.HandleUserCreated,
	})

	natsPubSubConsumer.Subscribe(natsconsumer.PubSubSubscriptionConfig{
		Subject: cfg.Nats.NatsSubjects.TaskCreatedEventSubject,
		Handler: natsHandler.HandleTaskCreated,
	})

	httpServer := httpserver.New(cfg.Server, statisticsUsecase)

	grpcServer := grpcserver.New(cfg.Server.GRPCServer, statisticsUsecase)

	app := &App{
		httpServer:         httpServer,
		natsPubSubConsumer: natsPubSubConsumer,
		grpcServer:         grpcServer,
	}

	return app, nil
}

func (a *App) Close(ctx context.Context) {
	err := a.httpServer.Stop()
	if err != nil {
		log.Println("failed to shutdown gRPC service", err)
	}

	err = a.grpcServer.Stop(ctx)
	if err != nil {
		log.Println("failed to shutdown gRPC service", err)
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()
	a.grpcServer.Run(ctx, errCh)
	a.httpServer.Run(errCh)
	a.natsPubSubConsumer.Start(ctx, errCh)
	log.Println(fmt.Sprintf("service %v started", serviceName))

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun

	case s := <-shutdownCh:
		log.Println(fmt.Sprintf("received signal: %v. Running graceful shutdown...", s))

		a.Close(ctx)
		log.Println("graceful shutdown completed!")
	}

	return nil
}
