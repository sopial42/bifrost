package candles

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	candlesSVC "github.com/sopial42/bifrost/internal/services/candles"
	domain "github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	appErrors "github.com/sopial42/bifrost/pkg/errors"
)

type candlesHandler struct {
	candlesSVC candlesSVC.Service
}

func SetHandler(e *echo.Echo, service candlesSVC.Service) {
	p := &candlesHandler{
		candlesSVC: service,
	}

	apiV1 := e.Group("/api/v1")
	{
		apiV1.GET("/candles/surrounding-dates", p.getSurroundingDates)
		apiV1.GET("/candles", p.getCandles)
		apiV1.PATCH("/candles", p.updateCandles)
		apiV1.POST("/candles", p.createcandles)
	}
}

type CandlesInputRequest struct {
	Candles []domain.Candle `json:"candles"`
}

func (p *candlesHandler) createcandles(context echo.Context) error {
	input := new(CandlesInputRequest)
	if err := context.Bind(input); err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	if len(input.Candles) == 0 {
		return appErrors.NewInvalidInput("invalid input, empty candles", nil)
	}

	candles, err := p.candlesSVC.CreateCandles(context.Request().Context(), &input.Candles)
	if err != nil {
		return fmt.Errorf("unable to create candles: %w", err)
	}

	return context.JSON(http.StatusCreated, candles)
}

func (p *candlesHandler) getCandles(context echo.Context) error {
	ctx := context.Request().Context()

	pair := common.Pair(context.QueryParam("pair"))
	interval := common.Interval(context.QueryParam("interval"))

	if pair == "" || interval == "" {
		return appErrors.NewInvalidInput("invalid input, pair and interval are required", nil)
	}

	limit, err := strconv.Atoi(context.QueryParam("limit"))
	if err != nil {
		return appErrors.NewInvalidInput("invalid input, limit is required", err)
	}

	if limit <= 0 {
		return appErrors.NewInvalidInput("invalid input, limit must be greater than 0", nil)
	}

	startDate := context.QueryParam("start_date")
	if startDate == "" {
		return appErrors.NewInvalidInput("invalid input, start_date is required", nil)
	}

	startDateParsed, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return appErrors.NewInvalidInput("invalid input, start_date is required in RFC3339 format", err)
	}

	candles, hasMore, err := p.candlesSVC.GetCandles(ctx, pair, interval, &startDateParsed, limit)
	if err != nil {
		return fmt.Errorf("unable to get candles: %w", err)
	}

	return context.JSON(http.StatusOK, map[string]interface{}{
		"candles":  candles,
		"has_more": hasMore,
	})
}

func (p *candlesHandler) getSurroundingDates(context echo.Context) error {
	pair := common.Pair(context.QueryParam("pair"))
	interval := common.Interval(context.QueryParam("interval"))

	if pair == "" || interval == "" {
		return appErrors.NewInvalidInput("invalid input, pair and interval are required", nil)
	}

	firstDate, lastDate, err := p.candlesSVC.GetSurroundingDates(context.Request().Context(), pair, interval)
	if err != nil {
		return fmt.Errorf("unable to get first and last dates: %w", err)
	}

	return context.JSON(http.StatusOK, map[string]string{
		"first_date": firstDate.String(),
		"last_date":  lastDate.String(),
	})
}

func (p *candlesHandler) updateCandles(context echo.Context) error {
	input := new(CandlesInputRequest)
	if err := context.Bind(input); err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	if len(input.Candles) == 0 {
		return appErrors.NewInvalidInput("invalid input, empty candles", nil)
	}

	candles, err := p.candlesSVC.UpdateCandles(context.Request().Context(), &input.Candles)
	if err != nil {
		return fmt.Errorf("unable to update candles: %w", err)
	}

	return context.JSON(http.StatusOK, candles)
}
