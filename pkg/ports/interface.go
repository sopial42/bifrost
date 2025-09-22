package ports

import (
	"context"
	"time"

	bsDomain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	"github.com/sopial42/bifrost/pkg/domains/positions"
)

type Client interface {
	Candles
	BuySignals
	Positions
}

type PriceRequest map[common.Pair][]candles.Date
type PriceResponse map[common.Pair]map[PriceRequestDate]float64

type PriceRequestDate string

type Candles interface {
	// CreateCandles insert candles in the database, if a candle already exists, it will be ignored
	// It returns only the newly inserted candles
	// The candle list is chunked by specified size or defaultChunckSize if set to <= 0
	CreateCandles(ctx context.Context, candles *[]candles.Candle, chunckSize int) (*[]candles.Candle, error)

	// GetCandles returns candles for a given pair and interval
	// Use startDate and endDate to filter candles by date
	// If limit param is > 0 use it as max candle count ot return, else use default sdk limit value
	// Return hasMore = true if there are more candles to fetch using nextCursor
	// Return nextCursor = the last candle date if there are more candles to fetch
	GetCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time, limit uint) (res *[]candles.Candle, hasMore bool, nextCursor *time.Time, err error)
	GetCandleByDate(ctx context.Context, pair common.Pair, interval common.Interval, date candles.Date) (res *candles.Candle, err error)
	GetCandlesMinuteClosePriceByDate(ctx context.Context, prices PriceRequest) (*PriceResponse, error)
	// GetCandlesByLastDate reverse the cursor, the next_cursor has to be used as last_date argument
	GetCandlesByLastDate(ctx context.Context, pair common.Pair, interval common.Interval, lastDate candles.Date, limit uint) (res *[]candles.Candle, hasMore bool, nextCursor *time.Time, err error)
	// GetCandlesByDate returns candles for a given pair and interval and date
	// UpdateCandleListRSI updates only the RSI for a list of candles
	// It returns the updated candles
	UpdateCandleListRSI(ctx context.Context, candles *[]candles.Candle) (*[]candles.Candle, error)

	// QuerySurroundingDates returns the first and last candle date for a given pair and interval
	// It returns 404 not found if no candles are found for the given pair and interval
	QuerySurroundingDates(ctx context.Context, pair common.Pair, interval common.Interval) (*candles.Date, *candles.Date, error)
}

type BuySignals interface {
	CreateBuySignals(ctx context.Context, buySignal *[]bsDomain.Details) (*[]bsDomain.Details, error)
	GetBuySignals(context.Context, common.Pair, common.Interval, bsDomain.Name, *time.Time) (res *[]bsDomain.Details, hasMore bool, nextCursor *time.Time, err error)
}

type Positions interface {
	CreatePositions(ctx context.Context, positions *[]positions.Details, chunckSize int) (*[]positions.Details, error)
}
