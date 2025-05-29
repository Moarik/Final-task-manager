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
	"taskManager/task/config"
	"taskManager/task/internal/adapter/http/service/handler"
	"taskManager/task/internal/usecase"
	"time"
)

const serverIPAddress = "0.0.0.0:%d"

type API struct {
	server *gin.Engine
	cfg    config.HTTPServer
	addr   string

	taskHandler *handler.Task
}

func New(cfg config.Server, taskUsecase usecase.Task) *API {
	gin.SetMode(cfg.HTTPServer.Mode)
	// Creating a new Gin Engine
	server := gin.New()

	// Applying middleware
	server.Use(gin.Recovery())

	taskHandler := handler.New(&taskUsecase)

	api := &API{
		server: server,
		cfg:    cfg.HTTPServer,
		addr:   fmt.Sprintf(serverIPAddress, cfg.HTTPServer.Port),

		taskHandler: taskHandler,
	}

	api.setupRoutes()

	return api
}

func (a *API) setupRoutes() {
	task := a.server.Group("/task")
	{
		task.GET("/user/all/:id", a.taskHandler.GetAllUserTasksByID)
		task.GET("/user/:id", a.taskHandler.GetUserTaskByID)
		task.POST("/user/", a.taskHandler.CreateUserTask)
		task.DELETE("/user/:id", a.taskHandler.DeleteUserTaskByID)
		task.PUT("/user/", a.taskHandler.UpdateUserTaskByID)
		task.GET("/:id", a.taskHandler.GetTaskByID)
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
