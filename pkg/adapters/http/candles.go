package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sopial42/bifrost/pkg/common/errors"
	"github.com/sopial42/bifrost/pkg/common/logger"
	"github.com/sopial42/bifrost/pkg/common/sdk"
	"github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	"github.com/sopial42/bifrost/pkg/ports"
)

const defaultCreateCandlesChunckSize = 5000
const defaultGetCandlesLimit = 5000

func (c *client) GetCandlesMinuteClosePriceByDate(ctx context.Context, prices ports.PriceRequest) (*ports.PriceResponse, error) {
	body, err := json.Marshal(prices)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal prices: %w", err)
	}

	res, err := c.Post(ctx, "/candles/minute-close-prices", body)
	if err != nil {
		return nil, fmt.Errorf("failed to get candles close price: %w", err)
	}

	postResponse := struct {
		Prices ports.PriceResponse `json:"prices"`
	}{}

	err = json.Unmarshal(res, &postResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal candles close price: %w", err)
	}
	return &postResponse.Prices, nil
}

func (c *client) CreateCandles(ctx context.Context, newCandles *[]candles.Candle, chunckSize int) (*[]candles.Candle, error) {
	if newCandles == nil || len(*newCandles) == 0 {
		return nil, nil
	}

	createdCandles := []candles.Candle{}
	// Chunk the candles list into specified size or 1000 if set to 0
	if chunckSize <= 0 {
		chunckSize = defaultCreateCandlesChunckSize
	}

	chuncks := sdk.CreateChunk(newCandles, chunckSize)
	if chuncks == nil {
		return nil, nil
	}

	for _, chunk := range *chuncks {
		input := map[string]interface{}{
			"candles": chunk,
		}
		body, err := json.Marshal(input)
		if err != nil {
			return nil, errors.NewUnexpected("failed to marshal candles", err)
		}

		res, err := c.Post(ctx, "/candles", body)
		if err != nil {
			return nil, err
		}

		postResponse := struct {
			Candles []candles.Candle `json:"candles"`
		}{}

		err = json.Unmarshal(res, &postResponse)
		if err != nil {
			return nil, errors.NewUnexpected("create failed to unmarshal candles while createChunck", err)
		}

		createdCandles = append(createdCandles, postResponse.Candles...)
	}

	return &createdCandles, nil
}

func (c *client) UpdateCandleListRSI(ctx context.Context, candlesRSIs *[]candles.Candle) (*[]candles.Candle, error) {
	if candlesRSIs == nil || len(*candlesRSIs) == 0 {
		return nil, nil
	}

	updatedCandles := []candles.Candle{}
	chuncks := sdk.CreateChunk(candlesRSIs, defaultCreateCandlesChunckSize)
	if chuncks == nil {
		return nil, nil
	}

	for _, chunk := range *chuncks {
		input := map[string]interface{}{
			"candles": chunk,
		}
		body, err := json.Marshal(input)
		if err != nil {
			return nil, errors.NewUnexpected("failed to marshal candles", err)
		}

		res, err := c.Patch(ctx, "/candles/rsi", body)
		if err != nil {
			return nil, err
		}

		patchResponse := struct {
			Candles []candles.Candle `json:"candles"`
		}{}

		err = json.Unmarshal(res, &patchResponse)
		if err != nil {
			logger.GetLogger(ctx).Errorf("update candle list RSI failed to unmarshal PATCH response: %v", err)
			return nil, errors.NewUnexpected("update failed to unmarshal while create a chunck of candles", err)
		}

		updatedCandles = append(updatedCandles, patchResponse.Candles...)
	}

	return &updatedCandles, nil
}

func (c *client) GetCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time, limit uint) (*[]candles.Candle, bool, *time.Time, error) {
	queryValues := url.Values{}

	if limit <= 0 {
		limit = defaultGetCandlesLimit
	}

	queryValues.Add("pair", string(pair))
	queryValues.Add("interval", string(interval))
	queryValues.Add("limit", strconv.Itoa(int(limit)))

	if startDate != nil {
		queryValues.Add("start_date", startDate.Format(time.RFC3339))
	}

	res, err := c.Get(ctx, "/candles?"+queryValues.Encode())
	if err != nil {
		return nil, false, nil, fmt.Errorf("failed to get candles: %w", err)
	}

	candlesResponse := struct {
		Candles    []candles.Candle `json:"candles"`
		HasMore    bool             `json:"has_more"`
		NextCursor *time.Time       `json:"next_cursor"`
	}{}

	err = json.Unmarshal(res, &candlesResponse)
	if err != nil {
		return nil, false, nil, errors.NewUnexpected("failed to unmarshal GetCandles response", err)
	}

	return &candlesResponse.Candles, candlesResponse.HasMore, candlesResponse.NextCursor, nil
}

func (c *client) GetCandleByDate(ctx context.Context, pair common.Pair, interval common.Interval, date candles.Date) (res *candles.Candle, err error) {
	queryValues := url.Values{}

	queryValues.Add("pair", string(pair))
	queryValues.Add("interval", string(interval))
	queryValues.Add("start_date", date.String())
	queryValues.Add("limit", "1")

	body, err := c.Get(ctx, "/candles?"+queryValues.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to get candle: %w", err)
	}

	candleResponse := struct {
		Candles []candles.Candle `json:"candles"`
	}{}

	err = json.Unmarshal(body, &candleResponse)
	if err != nil {
		return nil, errors.NewUnexpected("failed to unmarshal GetCandleByDate response", err)
	}

	if len(candleResponse.Candles) == 0 {
		return nil, nil
	}
	return &candleResponse.Candles[0], nil
}

func (c *client) GetCandlesByLastDate(ctx context.Context, pair common.Pair, interval common.Interval, lastDate candles.Date, limit uint) (*[]candles.Candle, bool, *time.Time, error) {
	queryValues := url.Values{}

	if limit <= 0 {
		limit = defaultGetCandlesLimit
	}

	queryValues.Add("pair", string(pair))
	queryValues.Add("interval", string(interval))
	queryValues.Add("last_date", lastDate.String())
	queryValues.Add("limit", strconv.Itoa(int(limit)))

	res, err := c.Get(ctx, "/candles/from-last-date?"+queryValues.Encode())
	if err != nil {
		return nil, false, nil, fmt.Errorf("failed to get candles: %w", err)
	}

	candlesResponse := struct {
		Candles    []candles.Candle `json:"candles"`
		HasMore    bool             `json:"has_more"`
		NextCursor *time.Time       `json:"next_cursor"`
	}{}

	err = json.Unmarshal(res, &candlesResponse)
	if err != nil {
		return nil, false, nil, errors.NewUnexpected("failed to unmarshal GetCandlesByLastDate response", err)
	}

	return &candlesResponse.Candles, candlesResponse.HasMore, candlesResponse.NextCursor, nil
}

func (c *client) QuerySurroundingDates(ctx context.Context, pair common.Pair, interval common.Interval) (*candles.Date, *candles.Date, error) {
	res, err := c.Get(ctx, "/candles/surrounding-dates?pair="+string(pair)+"&interval="+string(interval))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query surrounding dates: %w", err)
	}

	surroundingDates := struct {
		FirstDate *candles.Date `json:"first_date"`
		LastDate  *candles.Date `json:"last_date"`
	}{}

	err = json.Unmarshal(res, &surroundingDates)
	if err != nil {
		return nil, nil, errors.NewUnexpected("failed to unmarshal surrounding dates", err)
	}

	return surroundingDates.FirstDate, surroundingDates.LastDate, nil
}
