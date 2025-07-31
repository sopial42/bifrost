package candles

import (
	"context"

	domain "github.com/bifrost/pkg/domains/candles"
	"github.com/bifrost/pkg/domains/common"
)

type Service interface {
	CreateCandles(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
	GetSurroundingDates(context.Context, common.Pair, common.Interval) (*domain.Date, *domain.Date, error)
}

type Persistence interface {
	InsertCandles(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
	QuerySurroundingDates(context.Context, common.Pair, common.Interval) (*domain.Date, *domain.Date, error)
}
