package candles

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	domain "github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	"github.com/sopial42/bifrost/pkg/logger"
)

type CandleDAO struct {
	bun.BaseModel `bun:"table:candles"`

	ID       uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Date     time.Time
	Pair     string
	Interval string
	Open     float64
	Close    float64
	High     float64
	Low      float64
	RSI      *json.RawMessage `bun:"type:jsonb"`
}

func candlesToCandlesDAO(ctx context.Context, candles *[]domain.Candle, isUpdate bool) *[]CandleDAO {
	log := logger.GetLogger(ctx)
	if candles == nil {
		return nil
	}

	res := make([]CandleDAO, len(*candles))
	for i, c := range *candles {
		res[i] = CandleDAO{
			Date:     time.Time(c.Date),
			Pair:     c.Pair.String(),
			Interval: c.Interval.String(),
			Open:     c.Open,
			Close:    c.Close,
			High:     c.High,
			Low:      c.Low,
		}

		if isUpdate {
			if c.ID == nil {
				log.Warnf("unable to update candle as ID is nil: %+v", c)
				continue
			}

			res[i].ID = uuid.UUID(*c.ID)
		} else {
			res[i].ID = uuid.New()
		}

		if c.RSI != nil {
			msg, err := json.Marshal(*c.RSI)
			if err != nil {
				log.Warnf("unable to marshal candle RSI: %v", err)
				continue
			}

			var msgg json.RawMessage
			res[i].RSI = &msgg
			*res[i].RSI = msg
		}
	}

	return &res
}

func candlesDAOsToCandlesDetails(ctx context.Context, candlesDAO *[]CandleDAO) *[]domain.Candle {
	log := logger.GetLogger(ctx)
	if candlesDAO == nil {
		return nil
	}

	candles := make([]domain.Candle, len(*candlesDAO))
	for i, c := range *candlesDAO {
		if c.ID == uuid.Nil {
			log.Warnf("unable to convert candle as ID is nil: %+v", c)
			continue
		}

		id := domain.ID(c.ID)

		candles[i] = domain.Candle{
			ID:       &id,
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
			err := json.Unmarshal(*c.RSI, &rsi)
			if err == nil {
				candles[i].RSI = &rsi
			}
		}
	}

	return &candles
}
