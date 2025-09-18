package candles

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/uptrace/bun"

	"github.com/sopial42/bifrost/pkg/common/logger"
	domain "github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
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
	if candlesDAO == nil {
		return nil
	}

	candles := make([]domain.Candle, 0)
	for _, c := range *candlesDAO {
		if newCandle := candleDAOToCandleDetails(ctx, &c); newCandle != nil {
			candles = append(candles, *newCandle)
		}
	}

	return &candles
}

func candleDAOToCandleDetails(ctx context.Context, candleDAO *CandleDAO) *domain.Candle {
	if candleDAO == nil {
		return nil
	}

	if candleDAO.ID == uuid.Nil {
		log.Warnf("unable to convert candle as ID is nil: %+v", candleDAO)
		return nil
	}

	id := domain.ID(candleDAO.ID)

	candle := domain.Candle{
		ID:       &id,
		Date:     domain.Date(candleDAO.Date),
		Pair:     common.Pair(candleDAO.Pair),
		Interval: common.Interval(candleDAO.Interval),
		Open:     candleDAO.Open,
		Close:    candleDAO.Close,
		High:     candleDAO.High,
		Low:      candleDAO.Low,
	}

	if candleDAO.RSI != nil {
		rsi := domain.RSI{}
		err := json.Unmarshal(*candleDAO.RSI, &rsi)
		if err == nil {
			candle.RSI = &rsi
		}
	}

	return &candle
}
