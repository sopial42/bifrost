package positions

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	positionsSVC "github.com/bifrost/internal/services/positions"
	buysignals "github.com/bifrost/pkg/domains/buySignals"
	domain "github.com/bifrost/pkg/domains/positions"
	appErrors "github.com/bifrost/pkg/errors"
)

type positionsHandler struct {
	positionsSVC positionsSVC.Service
}

func SetHandler(e *echo.Echo, service positionsSVC.Service) {
	p := &positionsHandler{
		positionsSVC: service,
	}

	apiV1 := e.Group("/api/v1")
	{
		apiV1.POST("/positions", p.createPositions)
	}
}

type NewPositionInput struct {
	Positions []InputPositions
}

type InputPositions struct {
	Name        domain.Name     `json:"name"`
	Fullname    domain.Fullname `json:"fullname"`
	BuySignalID uuid.UUID       `json:"buy_signal_id"`
	TP          float64         `json:"tp"`
	SL          float64         `json:"sl"`
	Metadata    map[string]any  `json:"metadata"`
}

func (p *positionsHandler) createPositions(context echo.Context) error {
	newPositionInput := new(NewPositionInput)
	if err := context.Bind(newPositionInput); err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	if len(newPositionInput.Positions) == 0 {
		return appErrors.NewInvalidInput("invalid input, empty positions", nil)
	}

	newPositionsDetails := make([]domain.Details, len(newPositionInput.Positions))
	for i, pos := range newPositionInput.Positions {
		newPositionsDetails[i] = domain.Details{
			BuySignalID: buysignals.ID(pos.BuySignalID),
			Name:        pos.Name,
			Fullname:    pos.Fullname,
			TP:          pos.TP,
			SL:          pos.SL,
			Metadata:    pos.Metadata,
		}
	}

	positions, err := p.positionsSVC.CreatePositions(context.Request().Context(), &newPositionsDetails)
	if err != nil {
		return appErrors.NewUnexpected("unable to create positions", err)
	}

	return context.JSON(http.StatusCreated, positions)
}
