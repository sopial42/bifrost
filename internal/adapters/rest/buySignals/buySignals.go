package buysignals

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	appErrors "github.com/bifrost/internal/common/errors"
	domain "github.com/bifrost/internal/domains/buySignals"
	"github.com/bifrost/internal/domains/common"
	buySignalsSVC "github.com/bifrost/internal/services/buySignals"
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
		apiV1.POST("/buySignals", p.createBuySignals)
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
	Date       time.Time         `json:"date"`
	Price      float64           `json:"price"`
	Metadata   map[string]any    `json:"metadata"`
}

func (p *buySignalsHandler) createBuySignals(context echo.Context) error {
	newBuySignalInput := new(NewBuySignalInput)
	if err := context.Bind(newBuySignalInput); err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	newBuySignalsDetails := make([]domain.Details, len(newBuySignalInput.BuySignals))
	for i, bs := range newBuySignalInput.BuySignals {
		newBuySignalsDetails[i] = domain.Details{
			Name:       bs.Name,
			BusinessID: bs.BusinessID,
			Fullname:   bs.Fullname,
			Pair:       bs.Pair,
			Date:       bs.Date,
			Price:      bs.Price,
			Metadata:   bs.Metadata,
		}
	}

	buySignals, err := p.buySignalsSVC.CreateBuySignals(context.Request().Context(), &newBuySignalsDetails)
	if err != nil {
		return appErrors.NewUnexpected("unable to create buySignals", err)
	}

	return context.JSON(http.StatusCreated, buySignals)
}
