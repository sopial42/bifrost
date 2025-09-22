package inProcess

import (
	"context"
	"time"

	"github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	"github.com/sopial42/bifrost/pkg/ports"
)

func (c *inProcessClient) GetCandlesMinuteClosePriceByDate(ctx context.Context, prices ports.PriceRequest) (*ports.PriceResponse, error) {
	return nil, nil
}

func (c *inProcessClient) CreateCandles(ctx context.Context, newCandles *[]candles.Candle, chunckSize int) (*[]candles.Candle, error) {
	return nil, nil
}

func (c *inProcessClient) UpdateCandleListRSI(ctx context.Context, candlesRSIs *[]candles.Candle) (*[]candles.Candle, error) {
	return nil, nil
}

func (c *inProcessClient) GetCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time, limit uint) (*[]candles.Candle, bool, *time.Time, error) {
	return nil, false, nil, nil
}

func (c *inProcessClient) GetCandleByDate(ctx context.Context, pair common.Pair, interval common.Interval, date candles.Date) (res *candles.Candle, err error) {
	return nil, nil
}

func (c *inProcessClient) GetCandlesByLastDate(ctx context.Context, pair common.Pair, interval common.Interval, lastDate candles.Date, limit uint) (*[]candles.Candle, bool, *time.Time, error) {
	return nil, false, nil, nil
}

func (c *inProcessClient) QuerySurroundingDates(ctx context.Context, pair common.Pair, interval common.Interval) (*candles.Date, *candles.Date, error) {
	return nil, nil, nil
}
