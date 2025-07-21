package candles

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bifrost/internal/common/logger"
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	domain "github.com/bifrost/internal/domains/candles"
	"github.com/bifrost/internal/domains/common"
)

type CandleDAO struct {
	bun.BaseModel `bun:"table:candles"`

	ID       uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	SerialID int64     `bun:"serial_id,autoincrement"`
	Date     time.Time
	Pair     string
	Interval string
	Open     float64
	Close    float64
	High     float64
	Low      float64
	RSI      json.RawMessage `bun:"type:jsonb"`
}

func candlesToCandlesDAO(ctx context.Context, candles *[]domain.Candle, isUpdate bool) *[]CandleDAO {
	log := logger.GetLogger(ctx)
	if candles == nil {
		return nil
	}

	res := make([]CandleDAO, len(*candles))
	for i, c := range *candles {
		res[i] = CandleDAO{
			ID:       uuid.New(),
			Date:     time.Time(c.Date),
			Pair:     c.Pair.String(),
			Interval: c.Interval.String(),
			Open:     c.Open,
			Close:    c.Close,
			High:     c.High,
			Low:      c.Low,
		}

		if isUpdate {
			res[i].ID = uuid.UUID(c.ID)
		}

		if c.RSI != nil {
			msg, err := json.Marshal(*c.RSI)
			if err != nil {
				log.Warnf("unable to marshal candle RSI: %v", err)
				continue
			}

			res[i].RSI = msg
		}
	}

	return &res
}

func candlesDAOsToCandlesDetails(candlesDAO *[]CandleDAO) *[]domain.Candle {
	if candlesDAO == nil {
		return nil
	}

	candles := make([]domain.Candle, len(*candlesDAO))
	for i, c := range *candlesDAO {
		candles[i] = domain.Candle{
			ID:       domain.ID(c.ID),
			Date:     domain.Date(c.Date),
			Pair:     common.Pair(c.Pair),
			Interval: common.Interval(c.Interval),
			Open:     c.Open,
			Close:    c.Close,
			High:     c.High,
			Low:      c.Low,
		}

		if c.RSI != nil {
			rsi := domain.RSI{}
			err := json.Unmarshal(c.RSI, &rsi)
			if err == nil {
				candles[i].RSI = &rsi
			}
		}
	}

	return &candles
}
