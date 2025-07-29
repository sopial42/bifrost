package candles

import (
	"context"

	"github.com/uptrace/bun"

	domain "github.com/bifrost/internal/domains/candles"
	candlesSVC "github.com/bifrost/internal/services/candles"
	"github.com/bifrost/pkg/logger"
)

type pgPersistence struct {
	clientDB *bun.DB
}

func NewPersistence(client *bun.DB) candlesSVC.Persistence {
	return &pgPersistence{clientDB: client}
}

func (c *pgPersistence) InsertCandles(ctx context.Context, candles *[]domain.Candle) (*[]domain.Candle, error) {
	log := logger.GetLogger(ctx)

	if candles == nil {
		return &[]domain.Candle{}, nil
	}

	candlesDAO := candlesToCandlesDAO(ctx, candles, false)
	_, err := c.clientDB.NewInsert().Model(candlesDAO).Exec(ctx)
	if err != nil {
		log.Errorf("Query refused, throw pgError, %v", err)
		return &[]domain.Candle{}, err
	}

	log.Debugf("Insert done")
	return candlesDAOsToCandlesDetails(candlesDAO), nil
}
