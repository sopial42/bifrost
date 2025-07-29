package candles

import (
	"context"
	"fmt"

	domain "github.com/bifrost/internal/domains/candles"
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
