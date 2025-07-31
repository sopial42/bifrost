package candles

import (
	"fmt"
	"net/http"

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
		apiV1.POST("/candles", p.createcandles)
		apiV1.GET("/candles/surrounding-dates", p.getSurroundingDates)
		// apiV1.GET("/candles", p.getcandles)
	}
}

type NewCandleInput struct {
	Candles []domain.Candle `json:"candles"`
}

func (p *candlesHandler) createcandles(context echo.Context) error {
	input := new(NewCandleInput)
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
