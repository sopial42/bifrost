package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sopial42/bifrost/pkg/domains/candles"
	"github.com/sopial42/bifrost/pkg/domains/common"
	"github.com/sopial42/bifrost/pkg/errors"
	"github.com/sopial42/bifrost/pkg/logger"
)

const defaultCreateCandlesChunckSize = 1000
const defaultGetCandlesLimit = 100000

type Candles interface {
	// CreateCandles insert candles in the database, if a candle already exists, it will be ignored
	// It returns only the newly inserted candles
	// The candle list is chunked by specified size or defaultChunckSize if set to <= 0
	CreateCandles(ctx context.Context, candles *[]candles.Candle, chunckSize int) (*[]candles.Candle, error)

	// GetCandles returns candles for a given pair and interval
	// Use startDate and endDate to filter candles by date
	// If limit param is > 0 use it as max candle count ot return, else use default sdk limit value
	// Return hasMore = true if there are more candles to fetch using nextCursor
	// Return nextCursor = the last candle date if there are more candles to fetch
	GetCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time, limit uint) (res *[]candles.Candle, hasMore bool, nextCursor *time.Time, err error)
	GetCandleByDate(ctx context.Context, pair common.Pair, interval common.Interval, date candles.Date) (res *candles.Candle, err error)
	// GetCandlesByLastDate reverse the cursor, the next_cursor has to be used as last_date argument
	GetCandlesByLastDate(ctx context.Context, pair common.Pair, interval common.Interval, lastDate candles.Date, limit uint) (res *[]candles.Candle, hasMore bool, nextCursor *time.Time, err error)

	// GetCandlesByDate returns candles for a given pair and interval and date
	// UpdateCandleListRSI updates only the RSI for a list of candles
	// It returns the updated candles
	UpdateCandleListRSI(ctx context.Context, candles *[]candles.Candle) (*[]candles.Candle, error)

	// QuerySurroundingDates returns the first and last candle date for a given pair and interval
	// It returns 404 not found if no candles are found for the given pair and interval
	QuerySurroundingDates(ctx context.Context, pair common.Pair, interval common.Interval) (*candles.Date, *candles.Date, error)
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

	chuncks := createChunk(newCandles, chunckSize)
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
	chuncks := createChunk(candlesRSIs, defaultCreateCandlesChunckSize)
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
