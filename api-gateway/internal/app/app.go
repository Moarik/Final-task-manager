package app

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"syscall"
	"taskManager/api-gateway/config"
	httpserver "taskManager/api-gateway/internal/adapter/http/server"
	"taskManager/api-gateway/internal/adapter/http/server/handler"
)

const serviceName = "api-gateway"

var (
	userConnection       = "user:4000"
	taskConnection       = "task:4001"
	statisticsConnection = "statistics:4002"
)

type App struct {
	httpServer *httpserver.API
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Println(fmt.Sprintf("starting %v service", serviceName))

	// grpc user connection
	userConn, err := grpc.Dial(userConnection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user: %v", err)
	}

	// grpc task connection
	taskConn, err := grpc.Dial(taskConnection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to task: %v", err)
	}

	// grpc statistics connection
	statisticsConn, err := grpc.Dial(statisticsConnection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to task: %v", err)
	}

	clientHandler := handler.New(cfg.Server, userConn, taskConn, statisticsConn)

	// http server
	httpServer := httpserver.New(cfg.Server, clientHandler)

	app := &App{
		httpServer: httpServer,
	}

	return app, nil
}

func (a *App) Close(ctx context.Context) {
	err := a.httpServer.Stop()
	if err != nil {
		log.Println("failed to shutdown gRPC service", err)
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()
	a.httpServer.Run(errCh)

	log.Println(fmt.Sprintf("service %v started", serviceName))

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

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
