package candles

import (
	"context"
	"fmt"
	"time"

	domain "github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	appErrors "github.com/sopial42/bifrost/pkg/errors"
)

type candlesService struct {
	persistence Persistence
}

func NewCandlesService(persistence Persistence) Service {
	return &candlesService{
		persistence: persistence,
	}
}

func (p *candlesService) CreateCandles(ctx context.Context, candles *[]domain.Candle) (*[]domain.Candle, error) {
	candles, err := p.persistence.InsertCandles(ctx, candles)
	if err != nil {
		return &[]domain.Candle{}, fmt.Errorf("unable to insert candles: %w", err)
	}

	return candles, nil
}

func (p *candlesService) GetSurroundingDates(ctx context.Context, pair common.Pair, interval common.Interval) (*domain.Date, *domain.Date, error) {
	firstDate, lastDate, err := p.persistence.QuerySurroundingDates(ctx, pair, interval)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get surrounding dates: %w", err)
	}

	if firstDate == nil || lastDate == nil {
		return nil, nil, appErrors.NewNotFound("no candles found for this pair and interval")
	}
	return firstDate, lastDate, nil
}

func (p *candlesService) GetCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time, limit int) (*[]domain.Candle, bool, error) {
	candles, hasMore, err := p.persistence.QueryCandles(ctx, pair, interval, startDate, limit)
	if err != nil {
		return nil, false, fmt.Errorf("unable to get candles: %w", err)
	}

	return candles, hasMore, nil
}

func (p *candlesService) UpdateCandles(ctx context.Context, candles *[]domain.Candle) (*[]domain.Candle, error) {
	candles, err := p.persistence.UpdateCandles(ctx, candles)
	if err != nil {
		return &[]domain.Candle{}, fmt.Errorf("unable to update candles: %w", err)
	}

	return candles, nil
}
