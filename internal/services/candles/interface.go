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
	GetCandles(context.Context, common.Pair, common.Interval, *time.Time, *time.Time, int) (*[]domain.Candle, bool, *time.Time, error)
	// GetCandlesFromLastDate reverse the cursor, the next_cursor has to be used as last_date argument
	GetCandlesFromLastDate(context.Context, common.Pair, common.Interval, *time.Time, int) (candles *[]domain.Candle, hasMore bool, nextCursor *time.Time, err error)
	GetCandlesThatHitTPOrSL(ctx context.Context, pair common.Pair, buyDate domain.Date, tp float64, sl float64) (*domain.Candle, *domain.Candle, error)
	UpdateCandlesRSI(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
}

type Persistence interface {
	InsertCandles(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
	UpdateCandlesRSI(context.Context, *[]domain.Candle) (*[]domain.Candle, error)
	QueryCandles(context.Context, common.Pair, common.Interval, *time.Time, *time.Time, int) (*[]domain.Candle, bool, *time.Time, error)
	QueryCandlesFromLastDate(context.Context, common.Pair, common.Interval, *time.Time, int) (*[]domain.Candle, bool, *time.Time, error)
	QueryCandlesThatHitTPOrSL(context.Context, common.Pair, domain.Date, float64, float64) (*domain.Candle, *domain.Candle, error)
	QuerySurroundingDates(context.Context, common.Pair, common.Interval) (*domain.Date, *domain.Date, error)
}
