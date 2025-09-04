package positions

import (
	"context"
	"fmt"

	candlesSVC "github.com/sopial42/bifrost/internal/services/candles"
	"github.com/sopial42/bifrost/pkg/domains/candles"
	domain "github.com/sopial42/bifrost/pkg/domains/positions"
	"github.com/sopial42/bifrost/pkg/logger"
)

type positionsService struct {
	persistence Persistence
	candles     candlesSVC.Service
}

func NewPositionsService(persistence Persistence, candles candlesSVC.Service) Service {
	return &positionsService{
		persistence: persistence,
		candles:     candles,
	}
}

func (p *positionsService) CreatePositions(ctx context.Context, positions *[]domain.Details) (*[]domain.Details, error) {
	pos, err := p.persistence.InsertPositions(ctx, positions)
	if err != nil {
		return &[]domain.Details{}, err
	}

	return pos, nil
}

func (p *positionsService) ComputeAllPositions(ctx context.Context) (int, error) {
	var err error
	var positions *[]domain.Details
	hasMore := true
	nextCursor := new(int64)
	log := logger.GetLogger(ctx)
	updatedPositionsCount := 0

	for hasMore {
		positions, hasMore, nextCursor, err = p.persistence.GetPositionsWithNoRatio(ctx, nextCursor, 1000)
		if err != nil {
			return 0, err
		}

		if positions == nil || len(*positions) == 0 {
			return updatedPositionsCount, nil
		}

		positionsWithRatios := make([]domain.Details, len(*positions))
		for i, position := range *positions {
			// get candle that hit the TP or the SL
			if position.BuySignal == nil {
				log.Errorf("no buy signal found for position ID: %v", position.ID)
				positionsWithRatios[i] = position
				continue
			}

			tpCandle, slCandle, err := p.candles.GetCandlesThatHitTPOrSL(ctx, position.BuySignal.Pair, candles.Date(position.BuySignal.Date), position.TP, position.SL)
			if err != nil {
				return 0, fmt.Errorf("unable to get candles that hit the TP or the SL: %w", err)
			}

			ratio := 0.0
			if tpCandle == nil && slCandle == nil {
				positionsWithRatios[i] = position
				continue
			}

			if tpCandle != nil && slCandle != nil {
				if tpCandle.Date.Before(slCandle.Date) {
					ratio = position.TP / position.BuySignal.Price
					log.Infof("tp: %v, bySignalPrice: %v", position.TP, position.BuySignal.Price)
				} else {
					ratio = position.SL / position.BuySignal.Price
				}
			} else if tpCandle == nil {
				ratio = position.SL / position.BuySignal.Price
			} else if slCandle == nil {
				ratio = position.TP / position.BuySignal.Price
			}

			position.Ratio = &ratio
			positionsWithRatios[i] = position
		}

		updatedPositions, err := p.persistence.InsertRatios(ctx, &positionsWithRatios)
		if err != nil {
			return 0, fmt.Errorf("unable to insert ratios: %w", err)
		}
		if updatedPositions != nil {
			updatedPositionsCount += len(*updatedPositions)
		}
	}

	return updatedPositionsCount, nil
}
