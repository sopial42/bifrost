package candles

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"

	"github.com/sopial42/bifrost/pkg/common/logger"
	domain "github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	candlesSVC "github.com/sopial42/bifrost/pkg/services/candles"
)

const backtestInterval = common.Interval("1m")

type pgPersistence struct {
	clientDB *bun.DB
}

func NewPersistence(client *bun.DB) candlesSVC.Persistence {
	return &pgPersistence{clientDB: client}
}

// InsertCandles insert candles in the database, if a candle already exists, it will be ignored
// Returns only the newly inserted candles
func (c *pgPersistence) InsertCandles(ctx context.Context, candles *[]domain.Candle) (*[]domain.Candle, error) {
	log := logger.GetLogger(ctx)

	if candles == nil {
		return &[]domain.Candle{}, nil
	}

	candlesDAO := candlesToCandlesDAO(ctx, candles, false)
	_, err := c.clientDB.NewInsert().
		On("CONFLICT (date, interval, pair) DO NOTHING").
		Returning("*").
		Model(candlesDAO).
		Exec(ctx)
	if err != nil {
		log.Errorf("Query refused, throw pgError, %v", err)
		return &[]domain.Candle{}, err
	}

	log.Debugf("Insert done (%d)", len(*candlesDAO))
	return candlesDAOsToCandlesDetails(ctx, candlesDAO), nil
}

func (c *pgPersistence) QueryCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time, lastDate *time.Time, limit int) (*[]domain.Candle, bool, *time.Time, error) {
	result := []CandleDAO{}
	request := c.clientDB.NewSelect().Model(&result).
		Where("pair = ?", pair).
		Where("interval = ?", interval).
		OrderExpr("date ASC")

	if startDate != nil && !startDate.IsZero() {
		request.Where("date >= ?", startDate)
	}

	if lastDate != nil && !lastDate.IsZero() {
		request.Where("date <= ?", lastDate)
	}

	if limit > 0 {
		request.Limit(limit + 1)
	}

	err := request.Scan(ctx)
	if err != nil {
		return nil, false, nil, fmt.Errorf("unable to perform db query: %v", err)
	}

	if limit <= 0 {
		return candlesDAOsToCandlesDetails(ctx, &result), false, nil, nil
	}

	hasMore := len(result) > limit
	var nextCursor *time.Time
	if hasMore {
		nextCursor = &result[len(result)-1].Date
		result = result[:limit]
	}

	return candlesDAOsToCandlesDetails(ctx, &result), hasMore, nextCursor, nil
}

func (c *pgPersistence) QueryCandlesFromLastDate(ctx context.Context, pair common.Pair, interval common.Interval, lastDate *time.Time, limit int) (*[]domain.Candle, bool, *time.Time, error) {
	result := []CandleDAO{}
	request := c.clientDB.NewSelect().Model(&result).
		Where("pair = ?", pair).
		Where("interval = ?", interval).
		OrderExpr("date DESC")

	if lastDate != nil && !lastDate.IsZero() {
		request.Where("date <= ?", lastDate)
	} else {
		return nil, false, nil, fmt.Errorf("last_date is required")
	}

	if limit > 0 {
		request.Limit(limit + 1)
	}

	err := request.Scan(ctx)
	if err != nil {
		return nil, false, nil, fmt.Errorf("unable to perform db query: %w", err)
	}

	// reverse the result
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}

	if limit <= 0 {
		return candlesDAOsToCandlesDetails(ctx, &result), false, nil, nil
	}

	hasMore := len(result) > limit
	var nextCursor *time.Time
	if hasMore {
		nextCursor = &result[0].Date
		result = result[1:]
	}

	return candlesDAOsToCandlesDetails(ctx, &result), hasMore, nextCursor, nil
}

func (c *pgPersistence) QueryCandlesPriceByDate(ctx context.Context, pair common.Pair, date domain.Date) (float64, error) {
	candleRes := []CandleDAO{}
	err := c.clientDB.NewSelect().Model(&candleRes).
		Where("pair = ?", pair).
		Where("interval = ?", common.Interval("1m")).
		Where("date = ?", date).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return 0.0, fmt.Errorf("unable to perform db query: %w", err)
	}

	if len(candleRes) == 0 {
		return 0.0, nil
	}

	return candleRes[0].Close, nil
}

func (c *pgPersistence) QueryCandlesThatHitTPOrSL(ctx context.Context, pair common.Pair, buyDate domain.Date, tp float64, sl float64) (*domain.Candle, *domain.Candle, error) {
	res := make([]*domain.Candle, 2)

	for i, search := range []PriceHitSearch{{
		Condition: PriceHitConditionTP,
		Price:     tp,
	}, {
		Condition: PriceHitConditionSL,
		Price:     sl,
	}} {
		candle, err := c.searchPrice(ctx, pair, buyDate, search)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to perform db query: %w", err)
		}

		if candle == nil {
			continue
		}

		res[i] = candle
	}

	return res[0], res[1], nil
}

type PriceHitSearch struct {
	Price     float64
	Condition PriceHitCondition
}

type PriceHitCondition string

const (
	PriceHitConditionTP PriceHitCondition = "high >= ?"
	PriceHitConditionSL PriceHitCondition = "low <= ?"
)

func (c *pgPersistence) searchPrice(ctx context.Context, pair common.Pair, buyDate domain.Date, search PriceHitSearch) (*domain.Candle, error) {
	candleRes := []CandleDAO{}
	log := logger.GetLogger(ctx)

	// first ensure candles are present in DB with the backtest interval
	dateReq := c.clientDB.NewSelect().Model(&candleRes).
		Where("pair = ?", pair).
		Where("interval = ?", backtestInterval).
		Where("date >= ?", time.Time(buyDate)).
		Limit(1)

	err := dateReq.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to perform db query: %w", err)
	}

	if len(candleRes) == 0 {
		now := time.Now()
		buyDateTime := time.Time(buyDate)

		// If buyDate is today or yesterday, it's normal to not have data yet
		// as market data is usually available with 1 day delay
		if buyDateTime.Year() == now.Year() && buyDateTime.Month() == now.Month() {
			daysDiff := now.Day() - buyDateTime.Day()
			if daysDiff <= 1 {
				log.Debugf("no available data for backtest. Pair: %s, interval: %s, buyDate: %s", pair, backtestInterval, buyDate)
				return nil, nil
			}
		}
		errMsg := fmt.Sprintf("no candles found for the backtest. Pair: %s, interval: %s, buyDate: %s", pair, backtestInterval, buyDate.String())
		return nil, errors.New(errMsg)
	}

	// Then perform the search
	baseReq := c.clientDB.NewSelect().
		Where("pair = ?", pair).
		Where("interval = ?", backtestInterval).
		Where("date >= ?", time.Time(buyDate)).
		Where(string(search.Condition), search.Price).
		OrderExpr("date ASC").
		Limit(1)

	err = baseReq.Model(&candleRes).
		Scan(ctx, &candleRes)
	if err != nil {
		return nil, fmt.Errorf("unable to perform db query: %w", err)
	}

	if len(candleRes) == 0 {
		return nil, nil
	}

	return candleDAOToCandleDetails(ctx, &candleRes[0]), nil
}

func (c *pgPersistence) QuerySurroundingDates(ctx context.Context, pair common.Pair, interval common.Interval) (*domain.Date, *domain.Date, error) {
	var row struct {
		First *time.Time `bun:"first_date"`
		Last  *time.Time `bun:"last_date"`
	}

	query := `
        SELECT MIN(date) as first_date, MAX(date) as last_date
        FROM candles
        WHERE pair = ? AND interval = ?
    `
	err := c.clientDB.NewRaw(query, pair, interval).Scan(ctx, &row)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to query the first and last candle date: %w", err)
	}

	if row.First == nil || row.Last == nil {
		return nil, nil, nil
	}

	firstDate := domain.Date(*row.First)
	lastDate := domain.Date(*row.Last)
	return &firstDate, &lastDate, nil
}

func (c *pgPersistence) UpdateCandlesRSI(ctx context.Context, candles *[]domain.Candle) (*[]domain.Candle, error) {
	log := logger.GetLogger(ctx)

	if candles == nil {
		return &[]domain.Candle{}, nil
	}

	candlesDAO := candlesToCandlesDAO(ctx, candles, true)
	err := c.clientDB.NewUpdate().
		Model(candlesDAO).
		Column("rsi").
		Bulk().
		Returning("candle_dao.*").
		Scan(ctx)
	if err != nil {
		log.Errorf("Query refused, throw pgError, %v", err)
		return &[]domain.Candle{}, err
	}

	return candlesDAOsToCandlesDetails(ctx, candlesDAO), nil
}
