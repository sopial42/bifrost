package ratios

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	ratiosSVC "github.com/sopial42/bifrost/internal/services/ratios"
	"github.com/sopial42/bifrost/pkg/domains/positions"
	domain "github.com/sopial42/bifrost/pkg/domains/ratios"
	appErrors "github.com/sopial42/bifrost/pkg/errors"
)

type ratiosHandler struct {
	ratiosSVC ratiosSVC.Service
}

func SetHandler(e *echo.Echo, service ratiosSVC.Service) {
	p := &ratiosHandler{
		ratiosSVC: service,
	}

	apiV1 := e.Group("/api/v1")
	{
		apiV1.POST("/ratios", p.createRatios)
	}
}

type NewRatioInput struct {
	Ratios []InputRatio `json:"ratios"`
}

type InputRatio struct {
	PositionID uuid.UUID `json:"position_id"`
	Ratio      float64   `json:"ratio"`
	Date       time.Time `json:"date"`
}

func (p *ratiosHandler) createRatios(context echo.Context) error {
	newRatioInput := new(NewRatioInput)
	if err := context.Bind(newRatioInput); err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	if len(newRatioInput.Ratios) == 0 {
		return appErrors.NewInvalidInput("invalid input, empty ratios", nil)
	}

	newRatiosDetails := make([]domain.Ratio, len(newRatioInput.Ratios))
	for i, r := range newRatioInput.Ratios {
		newRatiosDetails[i] = domain.Ratio{
			PositionID: positions.ID(r.PositionID),
			Ratio:      r.Ratio,
			Date:       r.Date,
		}
	}

	ratios, err := p.ratiosSVC.CreateRatios(context.Request().Context(), &newRatiosDetails)
	if err != nil {
		return appErrors.NewUnexpected("unable to create ratios", err)
	}

	return context.JSON(http.StatusCreated, ratios)
}
