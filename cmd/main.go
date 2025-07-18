package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gommonLog "github.com/labstack/gommon/log"

	"github.com/bifrost/internal/adapters/persistence"
	"github.com/bifrost/internal/common/errors"
	"github.com/bifrost/internal/common/config"
	"github.com/bifrost/internal/common/pinger"
	"github.com/bifrost/internal/common/logger"
)

const pingRoute = "/ping"

func main() {
	config := config.Load()

	// Init external clients
	pgClient := persistence.NewPGClient(config.DB)
	// Configure echo engine
	engine := echo.New()

	// Custom logger
	log := logger.NewLogger(config.Logger)
	defer log.Sync() //nolint:errcheck
	logger.SetLoggerMiddlewareEcho(engine, log)
	logger.SetHTTPLoggerMiddlewareEcho(engine, urlSkipper)
	engine.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogLevel: gommonLog.ERROR,
	}))

	errors.SetCustomErrorHandler(engine)
	pinger.SetNewPingers(engine, pingRoute, pgClient /*, mailClient*/)
	corsConfig := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{config.Cors.AllowOrigin},
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
	})
	engine.Use(corsConfig)

	// Start the server and handle shutdown
	go func() {
		if err := engine.Start(":8080"); err != nil {
			fmt.Printf("Shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Printf("Shutting down the server gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := engine.Shutdown(ctx); err != nil {
		fmt.Printf("Unable to shutdown server gracefully: %v\n", err)
		return
	}
}

func urlSkipper(c echo.Context) bool {
	return c.Path() == pingRoute &&
		c.Response().Status >= http.StatusOK
}
