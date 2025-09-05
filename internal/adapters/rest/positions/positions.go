package positions

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	positionsSVC "github.com/sopial42/bifrost/internal/services/positions"
	buysignals "github.com/sopial42/bifrost/pkg/domains/buySignals"
	domain "github.com/sopial42/bifrost/pkg/domains/positions"
	appErrors "github.com/sopial42/bifrost/pkg/errors"
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
		apiV1.POST("/positions/compute", p.computeAllPositions)
		apiV1.POST("/positions/compute/:id", p.computePosition)
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
	Ratio       *domain.Ratio   `json:"ratio,omitempty"`
}

func (p *positionsHandler) computePosition(context echo.Context) error {
	id := context.Param("id")
	idParsed, err := uuid.Parse(id)
	if err != nil {
		return appErrors.NewInvalidInput("invalid input", err)
	}

	updatedPosition, err := p.positionsSVC.ComputeRatio(context.Request().Context(), domain.ID(idParsed))
	if err != nil {
		return appErrors.NewUnexpected("unable to compute position", err)
	}

	return context.JSON(http.StatusOK, map[string]interface{}{
		"position": updatedPosition,
	})
}

func (p *positionsHandler) computeAllPositions(context echo.Context) error {
	updatedPositionsCount, err := p.positionsSVC.ComputeAllRatios(context.Request().Context())
	if err != nil {
		return appErrors.NewUnexpected("unable to compute all positions", err)
	}

	return context.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("%d positions computed", updatedPositionsCount),
	})
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
			Ratio:       pos.Ratio,
		}
	}

	positions, err := p.positionsSVC.CreatePositions(context.Request().Context(), &newPositionsDetails)
	if err != nil {
		return appErrors.NewUnexpected("unable to create positions", err)
	}

	return context.JSON(http.StatusCreated, map[string]interface{}{
		"positions": positions,
	})
}
