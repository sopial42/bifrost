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

	persistence "github.com/sopial42/bifrost/pkg/adapters/persistence"

	HTTPHandler "github.com/sopial42/bifrost/internal/adapters/httpserver"

	buySignalsPersistence "github.com/sopial42/bifrost/pkg/adapters/persistence/buySignals"
	buySignalsSVC "github.com/sopial42/bifrost/pkg/services/buySignals"

	candlesPersistence "github.com/sopial42/bifrost/pkg/adapters/persistence/candles"
	candlesSVC "github.com/sopial42/bifrost/pkg/services/candles"

	positionsPersistence "github.com/sopial42/bifrost/pkg/adapters/persistence/positions"
	positionsSVC "github.com/sopial42/bifrost/pkg/services/positions"

	"github.com/sopial42/bifrost/pkg/common/config"
	"github.com/sopial42/bifrost/pkg/common/errors"
	"github.com/sopial42/bifrost/pkg/common/logger"
	"github.com/sopial42/bifrost/pkg/common/pinger"
)

const pingRoute = "/ping"

func main() {
	config := config.Load()

	// Init external clients
	pgClient := persistence.NewPGClient(config.DB)
	// Configure echo engine
	engine := echo.New()

	buySignalsPersistence := buySignalsPersistence.NewPersistence(pgClient.Client)
	buySignalsService := buySignalsSVC.NewBuySignalsService(buySignalsPersistence)

	candlesPersistence := candlesPersistence.NewPersistence(pgClient.Client)
	candlesService := candlesSVC.NewCandlesService(candlesPersistence)

	positionsPersistence := positionsPersistence.NewPersistence(pgClient.Client)
	positionsService := positionsSVC.NewPositionsService(positionsPersistence, candlesService, buySignalsService)

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

	HTTPHandler.SetBuySignalsHTTPHandler(engine, buySignalsService)
	HTTPHandler.SetCandlesHTTPHandler(engine, candlesService)
	HTTPHandler.SetPositionsHTTPHandler(engine, positionsService)

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
