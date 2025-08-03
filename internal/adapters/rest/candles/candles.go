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
		apiV1.PATCH("/candles_rsi", p.updateCandlesRSI)
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

type getCandleResponse struct {
	Candles    *[]domain.Candle `json:"candles"`
	HasMore    bool             `json:"has_more"`
	NextCursor *time.Time       `json:"next_cursor,omitempty"`
}

func (p *candlesHandler) getCandles(context echo.Context) error {
	ctx := context.Request().Context()
	pair := common.Pair(context.QueryParam("pair"))
	interval := common.Interval(context.QueryParam("interval"))

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

	var startDate *time.Time
	startDateArg := context.QueryParam("start_date")
	if startDateArg == "" {
		startDate = nil
	} else {
		startDateParsed, err := time.Parse(time.RFC3339, startDateArg)
		if err != nil {
			return appErrors.NewInvalidInput("invalid input, start_date is required in RFC3339 format", err)
		}
		startDate = &startDateParsed
	}

	candles, hasMore, nextCursor, err := p.candlesSVC.GetCandles(ctx, pair, interval, startDate, limit)
	if err != nil {
		return fmt.Errorf("unable to get candles: %w", err)
	}

	return context.JSON(http.StatusOK, getCandleResponse{
		Candles:    candles,
		HasMore:    hasMore,
		NextCursor: nextCursor,
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

	return context.JSON(http.StatusOK, candles)
}
