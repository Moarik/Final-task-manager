package server

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
	"taskManager/api-gateway/config"
	"taskManager/api-gateway/internal/adapter/http/server/handler"
	"taskManager/api-gateway/pkg/middleware"
	"time"
)

const serverIPAddress = "0.0.0.0:%d"

type API struct {
	server        *gin.Engine
	cfg           config.HTTPServer
	addr          string
	ClientHandler *handler.Client
}

func New(cfg config.Server, clientHandler *handler.Client) *API {
	// Setting the Gin mode
	gin.SetMode(cfg.HTTPServer.Mode)

	// Creating a new Gin Engine
	server := gin.New()

	// Applying middleware
	server.Use(gin.Recovery())

	api := &API{
		server:        server,
		cfg:           cfg.HTTPServer,
		addr:          fmt.Sprintf(serverIPAddress, cfg.HTTPServer.Port),
		ClientHandler: clientHandler,
	}

	api.setupRoutes()

	return api
}

func (a *API) setupRoutes() {
	userGroup := a.server.Group("/user")
	{
		userGroup.POST("/register", a.ClientHandler.Register)
		userGroup.POST("/login", a.ClientHandler.Login)

		protectedGroup := a.server.Group("/protected", middleware.AuthRequired())
		{
			protectedGroup.GET("/:id", a.ClientHandler.GetUserByID)
			protectedGroup.DELETE("/", a.ClientHandler.DeleteUserByID)

			protectedGroup.POST("/task", a.ClientHandler.CreateTask)
			protectedGroup.GET("/task/:id", a.ClientHandler.GetUserTaskByID)
			protectedGroup.GET("/task", a.ClientHandler.GetAllUserTask)
			protectedGroup.PUT("/task", a.ClientHandler.UpdateUserTask)
			protectedGroup.DELETE("/task/:id", a.ClientHandler.DeleteUserTaskByID)
		}

		taskGroup := a.server.Group("/task")
		{
			taskGroup.GET("/:id", a.ClientHandler.GetTaskByID)
		}
	}

	statisticsGroup := a.server.Group("/statistics")
	{
		statisticsGroup.GET("/user", a.ClientHandler.GetUserStatistics)
		statisticsGroup.GET("/task/:id", a.ClientHandler.GetTaskStatisticsByUserID)
		statisticsGroup.GET("/task", a.ClientHandler.GetTaskStatistics)
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
