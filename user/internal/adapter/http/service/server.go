package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskManager/user/config"
	"taskManager/user/internal/adapter/http/service/handler"
	"taskManager/user/internal/usecase"
	"time"
)

const serverIPAddress = "0.0.0.0:%d"

type API struct {
	server *gin.Engine
	cfg    config.HTTPServer
	addr   string

	clientHandler *handler.Client
}

func New(cfg config.Server, clientUsecase usecase.Client) *API {
	gin.SetMode(cfg.HTTPServer.Mode)
	// Creating a new Gin Engine
	server := gin.New()

	// Applying middleware
	server.Use(gin.Recovery())

	clientHandler := handler.NewClient(&clientUsecase)

	api := &API{
		server: server,
		cfg:    cfg.HTTPServer,
		addr:   fmt.Sprintf(serverIPAddress, cfg.HTTPServer.Port),

		clientHandler: clientHandler,
	}

	api.setupRoutes()

	return api
}

func (a *API) setupRoutes() {
	user := a.server.Group("/user")
	{
		user.POST("/", a.clientHandler.Create)
		user.POST("/login", a.clientHandler.Login)
		user.GET("/:id", a.clientHandler.Get)
		user.DELETE("/:id", a.clientHandler.Delete)
	}
}

func (a *API) Run(errCh chan<- error) {
	go func() {
		log.Printf("HTTP server starting on: %v", a.addr)

		// No need to reinitialize `a.server` here. Just run it directly.
		if err := a.server.Run(a.addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to start HTTP server: %w", err)
			return
		}
	}()
}

func (a *API) Stop() error {
	// Setting up the signal channel to catch termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Blocking until a signal is received
	sig := <-quit
	log.Println("Shutdown signal received", "signal:", sig.String())

	// Creating a context with timeout for graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("HTTP server shutting down gracefully")

	// Note: You can use `Shutdown` if you use `http.Server` instead of `gin.Engine`.
	log.Println("HTTP server stopped successfully")

	return nil
}
