package candles

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	candlesSVC "github.com/bifrost/internal/services/candles"
	domain "github.com/bifrost/pkg/domains/candles"
	appErrors "github.com/bifrost/pkg/errors"
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

// func (p *candlesHandler) getcandles(context echo.Context) error {
// 	ctx := context.Request().Context()
// 	pair := context.QueryParam("pair")
// 	interval := context.QueryParam("interval")
// 	startDate := context.QueryParam("start_date")

// 	// ensure start_date is a valid date
// 	date, err := time.Parse(time.RFC3339, startDate)
// 	if err != nil {
// 		return appErrors.NewInvalidInput("invalid input, start_date is not a valid date", err)
// 	}

// 	if date.After(time.Now()) {
// 		return appErrors.NewInvalidInput("invalid input, start_date is in the future", nil)
// 	}

// 	candles, err := p.candlesSVC.GetCandles(ctx, pair, interval, date)
// 	if err != nil {
// 		return fmt.Errorf("unable to get candles: %w", err)

// 	}
// 	return context.JSON(http.StatusOK, "ok")
// }
