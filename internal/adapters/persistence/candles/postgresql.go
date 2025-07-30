package candles

import (
	"context"

	"github.com/uptrace/bun"

	candlesSVC "github.com/bifrost/internal/services/candles"
	domain "github.com/bifrost/pkg/domains/candles"
	"github.com/bifrost/pkg/logger"
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
