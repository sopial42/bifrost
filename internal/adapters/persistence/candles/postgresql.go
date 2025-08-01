package candles

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"

	candlesSVC "github.com/sopial42/bifrost/internal/services/candles"
	domain "github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	"github.com/sopial42/bifrost/pkg/logger"
)

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

func (c *pgPersistence) QueryCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time, limit int) (*[]domain.Candle, bool, error) {
	result := []CandleDAO{}
	request := c.clientDB.NewSelect().Model(&result).
		Where("pair = ?", pair).
		Where("interval = ?", interval).
		OrderExpr("date ASC")

	if startDate != nil && !startDate.IsZero() {
		request.Where("date >= ?", startDate)
	}

	err := request.Scan(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("unable to perform db query: %v", err)
	}

	return candlesDAOsToCandlesDetails(ctx, &result), false, nil
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
