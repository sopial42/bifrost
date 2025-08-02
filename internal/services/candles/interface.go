package candles

import (
	"context"
	"time"

	domain "github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
)

type Service interface {
	CreateCandles(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
	GetSurroundingDates(context.Context, common.Pair, common.Interval) (*domain.Date, *domain.Date, error)
	GetCandles(context.Context, common.Pair, common.Interval, *time.Time, int) (*[]domain.Candle, bool, error)
	UpdateCandles(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
}

type Persistence interface {
	InsertCandles(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
	UpdateCandles(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
	QueryCandles(context.Context, common.Pair, common.Interval, *time.Time, int) (*[]domain.Candle, bool, error)
	QuerySurroundingDates(context.Context, common.Pair, common.Interval) (*domain.Date, *domain.Date, error)
}
