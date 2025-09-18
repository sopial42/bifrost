package candles

import (
	"context"
	"fmt"
	"time"

	appErrors "github.com/sopial42/bifrost/pkg/common/errors"
	domain "github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
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

func (p *candlesService) GetCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time, lastDate *time.Time, limit int) (*[]domain.Candle, bool, *time.Time, error) {
	candles, hasMore, nextCursor, err := p.persistence.QueryCandles(ctx, pair, interval, startDate, lastDate, limit)
	if err != nil {
		return nil, false, nil, fmt.Errorf("unable to get candles: %w", err)
	}

	return candles, hasMore, nextCursor, nil
}

func (p *candlesService) GetCandlesFromLastDate(ctx context.Context, pair common.Pair, interval common.Interval, lastDate *time.Time, limit int) (*[]domain.Candle, bool, *time.Time, error) {
	candles, hasMore, nextCursor, err := p.persistence.QueryCandlesFromLastDate(ctx, pair, interval, lastDate, limit)
	if err != nil {
		return nil, false, nil, fmt.Errorf("unable to get candles: %w", err)
	}

	return candles, hasMore, nextCursor, nil
}

func (p *candlesService) GetCandlesThatHitTPOrSL(ctx context.Context, pair common.Pair, buyDate domain.Date, tp float64, sl float64) (*domain.Candle, *domain.Candle, error) {
	tpCandle, slCandle, err := p.persistence.QueryCandlesThatHitTPOrSL(ctx, pair, buyDate, tp, sl)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get candles that hit the TP or the SL: %w", err)
	}

	return tpCandle, slCandle, nil
}

func (p *candlesService) UpdateCandlesRSI(ctx context.Context, candles *[]domain.Candle) (*[]domain.Candle, error) {
	candles, err := p.persistence.UpdateCandlesRSI(ctx, candles)
	if err != nil {
		return &[]domain.Candle{}, fmt.Errorf("unable to update candles: %w", err)
	}

	return candles, nil
}

type PriceRequest map[common.Pair][]domain.Date
type PriceResponse map[common.Pair]map[PriceRequestDate]float64

type PriceRequestDate string

func (p *candlesService) GetCandlesMinuteClosePricesByDate(ctx context.Context, pricesRequest PriceRequest) (PriceResponse, error) {
	response := make(PriceResponse)
	for pair, dates := range pricesRequest {
		for _, date := range dates {
			newDate := common.Interval(common.M1).RoundDateToBeginingOfInterval(time.Time(date))
			price, err := p.persistence.QueryCandlesPriceByDate(ctx, pair, domain.Date(*newDate))
			if err != nil {
				return nil, fmt.Errorf("unable to get candles prices: %w", err)
			}

			if _, ok := response[pair]; !ok {
				response[pair] = make(map[PriceRequestDate]float64)
			}

			response[pair][PriceRequestDate(date.String())] = price
		}
	}

	return response, nil
}
