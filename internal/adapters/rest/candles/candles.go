package candles

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	domain "github.com/bifrost/internal/domains/candles"
	"github.com/bifrost/internal/domains/common"
	candlesSVC "github.com/bifrost/internal/services/candles"
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
	}
}

type NewCandleInput struct {
	Candles []Inputcandles `json:"candles"`
}

type Inputcandles struct {
	Pair     common.Pair     `json:"pair"`
	Interval common.Interval `json:"interval"`
	Date     time.Time       `json:"date"`
	Open     float64         `json:"open"`
	Close    float64         `json:"close"`
	High     float64         `json:"high"`
	Low      float64         `json:"low"`
}

type NewCandleInputC struct {
	Candles []domain.Candle `json:"candles"`
}

func (p *candlesHandler) createcandles(context echo.Context) error {
	newCandleInput := new(NewCandleInputC)
	if err := context.Bind(newCandleInput); err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	if len(newCandleInput.Candles) == 0 {
		return appErrors.NewInvalidInput("invalid input, empty candles", nil)
	}

	candles, err := p.candlesSVC.CreateCandles(context.Request().Context(), &newCandleInput.Candles)
	if err != nil {
		return appErrors.NewUnexpected("unable to create candles", err)
	}

	return context.JSON(http.StatusCreated, candles)
}
