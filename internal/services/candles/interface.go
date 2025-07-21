package candles

import (
	"context"

	domain "github.com/bifrost/internal/domains/candles"
)

type Service interface {
	CreateCandles(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
}

type Persistence interface {
	InsertCandles(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
}
