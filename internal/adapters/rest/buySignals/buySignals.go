package buysignals

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	buySignalsSVC "github.com/sopial42/bifrost/internal/services/buySignals"
	domain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/common"
	appErrors "github.com/sopial42/bifrost/pkg/errors"
)

type buySignalsHandler struct {
	buySignalsSVC buySignalsSVC.Service
}

func SetHandler(e *echo.Echo, service buySignalsSVC.Service) {
	p := &buySignalsHandler{
		buySignalsSVC: service,
	}

	apiV1 := e.Group("/api/v1")
	{
		apiV1.POST("/buy_signals", p.createBuySignals)
		apiV1.GET("/buy_signals", p.getBuySignals)
	}
}

type NewBuySignalInput struct {
	BuySignals []InputBuySignal `json:"buy_signals"`
}

type InputBuySignal struct {
	Name       domain.Name       `json:"name"`
	BusinessID domain.BusinessID `json:"business_id"`
	Fullname   domain.Fullname   `json:"fullname"`
	Pair       common.Pair       `json:"pair"`
	Interval   common.Interval   `json:"interval"`
	Date       time.Time         `json:"date"`
	Price      float64           `json:"price"`
	Metadata   map[string]any    `json:"metadata"`
}

func (p *buySignalsHandler) createBuySignals(context echo.Context) error {
	newBuySignalInput := new(NewBuySignalInput)
	if err := context.Bind(newBuySignalInput); err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	if len(newBuySignalInput.BuySignals) == 0 {
		return appErrors.NewInvalidInput("invalid input, empty buy signals", nil)
	}

	newBuySignalsDetails := make([]domain.Details, len(newBuySignalInput.BuySignals))
	for i, bs := range newBuySignalInput.BuySignals {
		newBuySignalsDetails[i] = domain.Details{
			Name:       bs.Name,
			BusinessID: bs.BusinessID,
			Fullname:   bs.Fullname,
			Pair:       bs.Pair,
			Interval:   bs.Interval,
			Date:       domain.Date(bs.Date),
			Price:      bs.Price,
			Metadata:   bs.Metadata,
		}
	}

	buySignals, err := p.buySignalsSVC.CreateBuySignals(context.Request().Context(), &newBuySignalsDetails)
	if err != nil && !errors.Is(err, appErrors.ErrAlreadyExists) {
		return fmt.Errorf("unable to create buySignals: %w", err)
	}

	return context.JSON(http.StatusCreated, map[string]any{
		"buy_signals": buySignals,
	})
}

func (p buySignalsHandler) getBuySignals(context echo.Context) (err error) {
	pair := context.QueryParam("pair")
	interval := context.QueryParam("interval")
	name := context.QueryParam("name")
	firstDate := context.QueryParam("first_date")
	limit := context.QueryParam("limit")

	pairParsed := common.Pair(pair)
	if !pairParsed.IsValid() {
		return appErrors.NewInvalidInput("invalid pair", nil)
	}

	intervalParsed := common.Interval(interval)
	if !intervalParsed.IsValid() {
		return appErrors.NewInvalidInput("invalid interval", nil)
	}

	firstDateParsed := time.Time{}
	if firstDate != "" {
		firstDateParsed, err = time.Parse(time.RFC3339, firstDate)
		if err != nil {
			return appErrors.NewInvalidInput("invalid first_date", err)
		}
	}

	limitParsed, err := strconv.Atoi(limit)
	if err != nil {
		return appErrors.NewInvalidInput("invalid limit", err)
	}

	nameParsed := domain.Name(name)
	if nameParsed == "" {
		return appErrors.NewInvalidInput("invalid name", nil)
	}

	buySignals, hasMore, nextCursor, err := p.buySignalsSVC.GetBuySignals(context.Request().Context(), pairParsed, intervalParsed, nameParsed, &firstDateParsed, limitParsed)
	if err != nil {
		return appErrors.NewUnexpected("unable to get buySignals", err)
	}

	return context.JSON(http.StatusOK, map[string]any{
		"buy_signals": buySignals,
		"has_more":    hasMore,
		"next_cursor": nextCursor,
	})
}
