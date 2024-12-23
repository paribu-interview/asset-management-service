package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/safayildirim/asset-management-service/internal/asset"
	"github.com/safayildirim/asset-management-service/internal/transaction"
	"github.com/safayildirim/asset-management-service/internal/transaction/scheduler"
	"github.com/safayildirim/asset-management-service/pkg/client/wallet"
	"github.com/safayildirim/asset-management-service/pkg/config"
	"github.com/safayildirim/asset-management-service/pkg/db"
	"github.com/safayildirim/asset-management-service/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Handler interface {
	RegisterRoutes(e *echo.Group)
}

type App struct {
	Config   config.Config
	DB       *gorm.DB
	Server   *echo.Echo
	Handlers []Handler
}

func New() *App {
	// Load configuration
	cfg := config.New()

	// Initialize GORM DB connection
	dbInstance, err := db.NewConnection(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	// Create Echo instance
	server := echo.New()

	server.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Logger.Error(err.Error(), zap.String("method", c.Request().Method),
			zap.String("path", c.Request().URL.Path))

		var httpError *echo.HTTPError
		errors.As(err, &httpError)

		// Default error handler behavior
		c.JSON(httpError.Code, err)
	}

	// Configure middleware
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	var handlers []Handler

	assetRepository := asset.NewRepository(dbInstance)
	walletClient := wallet.NewClient(cfg.WalletClient.BaseURL)
	assetService := asset.NewService(assetRepository, walletClient)
	assetHandler := asset.NewHandler(assetService)

	transactionRepository := transaction.NewRepository(dbInstance)
	transactionService := transaction.NewService(assetRepository, transactionRepository, walletClient)
	transactionHandler := transaction.NewHandler(transactionService)

	schedulerManager := scheduler.NewScheduler(cfg.Scheduler, assetService, transactionRepository)
	go schedulerManager.Start(context.Background())

	handlers = append(handlers, assetHandler, transactionHandler)

	return &App{Config: *cfg, DB: dbInstance, Server: server, Handlers: handlers}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	route := a.Server.Group("/api")

	for _, handler := range a.Handlers {
		handler.RegisterRoutes(route)
	}

	// Start the server in a goroutine
	go func() {
		port := fmt.Sprintf(":%d", a.Config.Http.Port)
		log.Logger.Info(fmt.Sprintf("Starting server on :%s", port))
		if err := a.Server.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("Could not start server: %v", err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Logger.Info("Shutting down server...")

	// Gracefully shut down the Echo server with a timeout
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	// Close database connection
	sqlDB, err := a.DB.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Logger.Info("Server exited gracefully")

	return nil
}
