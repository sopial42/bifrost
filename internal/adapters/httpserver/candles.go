package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	appErrors "github.com/sopial42/bifrost/pkg/common/errors"
	domain "github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	candlesSVC "github.com/sopial42/bifrost/pkg/services/candles"
)

type candlesHandler struct {
	candlesSVC candlesSVC.Service
}

func SetCandlesHTTPHandler(e *echo.Echo, service candlesSVC.Service) {
	p := &candlesHandler{
		candlesSVC: service,
	}

	apiV1 := e.Group("/api/v1")
	{
		apiV1.GET("/candles/surrounding-dates", p.getSurroundingDates)
		apiV1.GET("/candles", p.getCandles)
		apiV1.POST("/candles/minute-close-prices", p.getCandlesMinuteClosePricesByDate)
		apiV1.GET("/candles/from-last-date", p.getCandlesFromLastDate)
		apiV1.POST("/candles", p.createcandles)
		apiV1.PATCH("/candles/rsi", p.updateCandlesRSI)
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
	if err != nil && !errors.Is(err, appErrors.ErrAlreadyExists) {
		return fmt.Errorf("unable to create candles: %w", err)
	}

	return context.JSON(http.StatusCreated, map[string]interface{}{
		"candles": candles,
	})
}

func (p *candlesHandler) getCandles(context echo.Context) error {
	ctx := context.Request().Context()
	pair := common.Pair(context.QueryParam("pair"))
	interval := common.Interval(context.QueryParam("interval"))
	var startDate *time.Time
	startDateArg := context.QueryParam("start_date")
	var lastDate *time.Time
	lastDateArg := context.QueryParam("last_date")

	if pair == "" || interval == "" {
		return appErrors.NewInvalidInput("invalid input, pair and interval are required", nil)
	}

	var limit int
	var err error
	limitParam := context.QueryParam("limit")
	if limitParam == "" {
		limit = 0
	} else {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return appErrors.NewInvalidInput("invalid limit input, must be an integer", err)
		}
	}

	if startDateArg == "" {
		startDate = nil
	} else {
		startDateParsed, err := time.Parse(time.RFC3339, startDateArg)
		if err != nil {
			return appErrors.NewInvalidInput("invalid input, start_date is required in RFC3339 format", err)
		}
		startDate = &startDateParsed
	}

	if lastDateArg == "" {
		lastDate = nil
	} else {
		lastDateParsed, err := time.Parse(time.RFC3339, lastDateArg)
		if err != nil {
			return appErrors.NewInvalidInput("invalid input, last_date is required in RFC3339 format", err)
		}

		lastDate = &lastDateParsed
	}

	candles, hasMore, nextCursor, err := p.candlesSVC.GetCandles(ctx, pair, interval, startDate, lastDate, limit)
	if err != nil {
		return fmt.Errorf("unable to get candles: %w", err)
	}

	return context.JSON(http.StatusOK, map[string]any{
		"candles":     candles,
		"has_more":    hasMore,
		"next_cursor": nextCursor,
	})
}

// GetCandlesByLastDate reverse the cursor, the next_cursor has to be used as last_date argument
func (p *candlesHandler) getCandlesFromLastDate(context echo.Context) error {
	ctx := context.Request().Context()
	pair := common.Pair(context.QueryParam("pair"))
	interval := common.Interval(context.QueryParam("interval"))
	lastDateArg := context.QueryParam("last_date")

	if pair == "" || interval == "" {
		return appErrors.NewInvalidInput("invalid input, pair and interval are required", nil)
	}

	var lastDate *time.Time
	var limit int
	var err error
	limitParam := context.QueryParam("limit")
	if limitParam == "" {
		limit = 0
	} else {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return appErrors.NewInvalidInput("invalid limit input, must be an integer", err)
		}
	}

	if lastDateArg == "" {
		return appErrors.NewInvalidInput("invalid input, last_date is required", nil)
	} else {
		lastDateParsed, err := time.Parse(time.RFC3339, lastDateArg)
		if err != nil {
			return appErrors.NewInvalidInput("invalid input, last_date is required in RFC3339 format", err)
		}

		lastDate = &lastDateParsed
	}

	candles, hasMore, nextCursor, err := p.candlesSVC.GetCandlesFromLastDate(ctx, pair, interval, lastDate, limit)
	if err != nil {
		return fmt.Errorf("unable to get candles: %w", err)
	}

	return context.JSON(http.StatusOK, map[string]any{
		"candles":     candles,
		"has_more":    hasMore,
		"next_cursor": nextCursor,
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

	return context.JSON(http.StatusOK, map[string]any{
		"first_date": firstDate,
		"last_date":  lastDate,
	})
}

func (p *candlesHandler) updateCandlesRSI(context echo.Context) error {
	ctx := context.Request().Context()
	input := new(CandlesInputRequest)

	if err := context.Bind(input); err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	if len(input.Candles) == 0 {
		return appErrors.NewInvalidInput("invalid input, empty candles", nil)
	}

	candles, err := p.candlesSVC.UpdateCandlesRSI(ctx, &input.Candles)
	if err != nil {
		return fmt.Errorf("unable to update candles: %w", err)
	}

	return context.JSON(http.StatusOK, map[string]interface{}{
		"candles": candles,
	})
}

// Return the closing price of the candle for a given pair and date using 1 minute interval data
func (p *candlesHandler) getCandlesMinuteClosePricesByDate(context echo.Context) error {
	ctx := context.Request().Context()

	input := new(candlesSVC.PriceRequest)
	if err := context.Bind(input); err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	candlesPrices, err := p.candlesSVC.GetCandlesMinuteClosePricesByDate(ctx, *input)
	if err != nil {
		return fmt.Errorf("unable to get candles prices: %w", err)
	}

	return context.JSON(http.StatusOK, map[string]interface{}{
		"prices": candlesPrices,
	})
}
