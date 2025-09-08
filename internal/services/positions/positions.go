package positions

import (
	"context"
	"fmt"
	"time"

	buySignalsSVC "github.com/sopial42/bifrost/internal/services/buySignals"
	candlesSVC "github.com/sopial42/bifrost/internal/services/candles"
	"github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	domain "github.com/sopial42/bifrost/pkg/domains/positions"
	"github.com/sopial42/bifrost/pkg/logger"
)

type positionsService struct {
	persistence Persistence
	candles     candlesSVC.Service
	buySignals  buySignalsSVC.Service
}

func NewPositionsService(persistence Persistence, candles candlesSVC.Service, buySignals buySignalsSVC.Service) Service {
	return &positionsService{
		persistence: persistence,
		candles:     candles,
		buySignals:  buySignals,
	}
}

func (p *positionsService) CreatePositions(ctx context.Context, positions *[]domain.Details) (*[]domain.Details, error) {
	pos, err := p.persistence.InsertPositions(ctx, positions)
	if err != nil {
		return &[]domain.Details{}, err
	}

	return pos, nil
}

func (p *positionsService) ComputeRatio(ctx context.Context, id domain.ID) (*domain.Details, error) {
	position, err := p.persistence.GetPositionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get position by ID: %w", err)
	}

	computePosition, err := p.computeRatio(ctx, position)
	if err != nil {
		return nil, fmt.Errorf("unable to compute position: %w", err)
	}

	position.Ratio = computePosition
	return position, nil
}

func (p *positionsService) ComputeAllRatios(ctx context.Context) (int, error) {
	var err error
	var positions *[]domain.Details
	hasMore := true
	nextCursor := new(int64)
	log := logger.GetLogger(ctx)
	updatedPositionsCount := 0

	for hasMore {
		positions, hasMore, nextCursor, err = p.persistence.GetPositionsWithNoRatio(ctx, nextCursor, 100)
		if err != nil {
			return 0, err
		}

		if positions == nil || len(*positions) == 0 {
			return updatedPositionsCount, nil
		}

		log.Infof("Start compute positions: %v", len(*positions))
		positionsWithRatios := make([]domain.Details, 0)
		for _, position := range *positions {
			// get candle that hit the TP or the SL
			if position.BuySignal == nil {
				log.Errorf("no buy signal found for position ID: %v", position.ID)
				continue
			}

			ratio, err := p.computeRatio(ctx, &position)
			if err != nil {
				return 0, fmt.Errorf("unable to compute position: %w", err)
			}

			if ratio == nil || ratio.Value == 0 {
				continue
			}

			position.Ratio = ratio
			positionsWithRatios = append(positionsWithRatios, position)
		}

		updatedPositions, err := p.persistence.InsertRatios(ctx, &positionsWithRatios)
		if err != nil {
			return 0, fmt.Errorf("unable to insert ratios: %w", err)
		}

		if updatedPositions != nil {
			log.Infof("Updated positions: %v", len(*updatedPositions))
			updatedPositionsCount += len(*updatedPositions)
		}
	}

	return updatedPositionsCount, nil
}

func (p *positionsService) computeRatio(ctx context.Context, position *domain.Details) (*domain.Ratio, error) {
	result := domain.Ratio{}
	log := logger.GetLogger(ctx)

	if position.BuySignal == nil {
		return nil, fmt.Errorf("buy signal is required")
	}

	buyDate := common.AddOneInterval(time.Time(position.BuySignal.Date), position.BuySignal.Interval)
	if buyDate == nil {
		return nil, fmt.Errorf("unable to add one interval to the buy date, interval: %s", position.BuySignal.Interval)
	}

	tpCandle, slCandle, err := p.candles.GetCandlesThatHitTPOrSL(ctx, position.BuySignal.Pair, candles.Date(*buyDate), position.TP, position.SL)
	if err != nil {
		return nil, fmt.Errorf("unable to get candles that hit the TP or the SL: %w", err)
	}

	if tpCandle == nil && slCandle == nil {
		log.Debugf("no candles that hit. position: %+v, bs: %+v", position, position.BuySignal)
		return nil, nil
	}

	if tpCandle != nil && slCandle != nil {
		if tpCandle.Date.Before(slCandle.Date) {
			result.Value = position.TP / position.BuySignal.Price
			result.Date = tpCandle.Date
		} else {
			result.Value = position.SL / position.BuySignal.Price
			result.Date = slCandle.Date
		}
	} else if tpCandle == nil {
		result.Value = position.SL / position.BuySignal.Price
		result.Date = slCandle.Date
	} else if slCandle == nil {
		result.Value = position.TP / position.BuySignal.Price
		result.Date = tpCandle.Date
	}

	return &result, nil
}

func (p *positionsService) CreatePositionsWithBuySignals(ctx context.Context, positions *[]domain.Details) (*[]domain.Details, error) {
	addedPositions := make([]domain.Details, 0)
	for _, position := range *positions {
		if position.BuySignal == nil {
			return nil, fmt.Errorf("buy signal is required")
		}

		newBS, err := p.buySignals.UpsertBuySignals(ctx, *position.BuySignal)
		if err != nil {
			return nil, fmt.Errorf("unable to create buy signal: %w", err)
		}

		if newBS == nil || len(*newBS) == 0 {
			return nil, fmt.Errorf("unable to upsert buy signal")
		}

		currentBS := (*newBS)[0]
		position.BuySignalID = *currentBS.ID
		position.BuySignal.ID = currentBS.ID
		var newPos *domain.Details

		newPos, err = p.persistence.UpsertPosition(ctx, &position)
		if err != nil {
			return nil, fmt.Errorf("unable to upsert position: %w", err)
		}

		if newPos == nil {
			return nil, fmt.Errorf("unable to upsert position")
		}

		newPos.BuySignal = &currentBS
		ratio, err := p.computeRatio(ctx, newPos)
		if err != nil {
			return nil, fmt.Errorf("unable to compute ratio: %w", err)
		}
		newPos.Ratio = ratio
		addedPositions = append(addedPositions, *newPos)
	}

	return &addedPositions, nil
}
