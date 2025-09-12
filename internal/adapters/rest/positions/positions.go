package positions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	positionsSVC "github.com/sopial42/bifrost/internal/services/positions"
	buysignals "github.com/sopial42/bifrost/pkg/domains/buySignals"
	domain "github.com/sopial42/bifrost/pkg/domains/positions"
	appErrors "github.com/sopial42/bifrost/pkg/errors"
	"github.com/sopial42/bifrost/pkg/logger"
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
		apiV1.POST("/positions/compute/with-buy-signals", p.createPositionsWithBuySignals)
		apiV1.POST("/positions/compute/all", p.computeAllPositions)
		apiV1.POST("/positions/compute/:id", p.computePosition)
	}
}

type NewPositionInput struct {
	Positions []InputPositions
}

type InputPositions struct {
	Name         domain.Name          `json:"name"`
	Fullname     domain.Fullname      `json:"fullname"`
	BuySignalID  uuid.UUID            `json:"buy_signal_id"`
	TP           float64              `json:"tp"`
	SL           float64              `json:"sl"`
	Metadata     map[string]any       `json:"metadata"`
	Ratio        *domain.Ratio        `json:"ratio,omitempty"`
	WinlossRatio *domain.WinLossRatio `json:"winloss_ratio,omitempty"`
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
			BuySignalID:  buysignals.ID(pos.BuySignalID),
			Name:         pos.Name,
			Fullname:     pos.Fullname,
			TP:           pos.TP,
			SL:           pos.SL,
			Metadata:     pos.Metadata,
			Ratio:        pos.Ratio,
			WinlossRatio: pos.WinlossRatio,
		}
	}

	positions, err := p.positionsSVC.CreatePositions(context.Request().Context(), &newPositionsDetails)
	if err != nil && !errors.Is(err, appErrors.ErrAlreadyExists) {
		return appErrors.NewUnexpected("unable to create positions", err)
	}

	return context.JSON(http.StatusCreated, map[string]interface{}{
		"positions": positions,
	})
}

type NewPositionInputWithBuySignals struct {
	Positions []domain.Details `json:"positions"`
}

func (p *positionsHandler) createPositionsWithBuySignals(context echo.Context) error {
	input := new(NewPositionInputWithBuySignals)
	log := logger.GetLogger(context.Request().Context())

	if err := context.Bind(input); err != nil {
		log.Debugf("invalid input: %v", err)
		return appErrors.NewInvalidInput("invalid input", err)
	}

	if len(input.Positions) == 0 {
		log.Debugf("invalid input, empty positions")
		return appErrors.NewInvalidInput("invalid input, empty positions", nil)
	}

	for _, pos := range input.Positions {
		var err error
		if pos.ID != nil {
			err = errors.Join(err, appErrors.NewInvalidInput("position.id must not be provided", nil))
		}

		if pos.Fullname == "" {
			err = errors.Join(err, appErrors.NewInvalidInput("position.fullname is required", nil))
		}

		if pos.Name == "" {
			err = errors.Join(err, appErrors.NewInvalidInput("position.name is required", nil))
		}

		if pos.TP == 0 {
			err = errors.Join(err, appErrors.NewInvalidInput("position.tp is required", nil))
		}

		if pos.SL == 0 {
			err = errors.Join(err, appErrors.NewInvalidInput("position.sl is required", nil))
		}

		if pos.TP <= pos.SL || pos.SL <= 0 || pos.TP <= 0 {
			err = errors.Join(err, appErrors.NewInvalidInput("position.TP or position.SL is invalid: It should be greater than 0. position.TP should be greater than position.SL.", nil))
		}

		if pos.BuySignal == nil {
			err = errors.Join(err, appErrors.NewInvalidInput("position.buy_signal is required", nil))
			return appErrors.NewInvalidInput(fmt.Sprintf("invalid input : %v", err), err)
		}

		if pos.BuySignal.ID != nil {
			err = errors.Join(err, appErrors.NewInvalidInput("buy_signal.id must not be provided", nil))
		}

		if pos.BuySignal.BusinessID == "" {
			err = errors.Join(err, appErrors.NewInvalidInput("buy_signal.business_id should be provided", nil))
		}

		if pos.BuySignal.Fullname == "" {
			err = errors.Join(err, appErrors.NewInvalidInput("position.buy_signal.fullname is required", nil))
		}

		if pos.BuySignal.Name == "" {
			err = errors.Join(err, appErrors.NewInvalidInput("position.buy_signal.name is required", nil))
		}

		if pos.BuySignal.Pair == "" {
			err = errors.Join(err, appErrors.NewInvalidInput("position.buy_signal.pair is required", nil))
		}

		if pos.BuySignal.Interval == "" {
			err = errors.Join(err, appErrors.NewInvalidInput("position.buy_signal.interval is required", nil))
		}

		if pos.BuySignal.Date == (buysignals.Date{}) {
			err = errors.Join(err, appErrors.NewInvalidInput("position.buy_signal.date is required", nil))
		}

		if pos.BuySignal.Price == 0 {
			err = errors.Join(err, appErrors.NewInvalidInput("position.buy_signal.price is required", nil))
		}

		if err != nil {
			return appErrors.NewInvalidInput(fmt.Sprintf("invalid input : %v", err), err)
		}
	}

	positions, err := p.positionsSVC.CreatePositionsWithBuySignals(context.Request().Context(), &input.Positions)
	if err != nil && !errors.Is(err, appErrors.ErrAlreadyExists) {
		return appErrors.NewUnexpected("unable to create positions with buy signals", err)
	}

	return context.JSON(http.StatusCreated, map[string]interface{}{
		"positions": positions,
	})
}
