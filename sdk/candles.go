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
	"github.com/sopial42/bifrost/pkg/sdk"
)

const defaultCreateCandlesChunckSize = 1000
const defaultGetCandlesLimit = 1000

type Candles interface {
	// CreateCandles insert candles in the database, if a candle already exists, it will be ignored
	// It returns only the newly inserted candles
	// The candle list is chunked by specified size or defaultChunckSize if set to <= 0
	CreateCandles(ctx context.Context, candles *[]candles.Candle, chunckSize int) (*[]candles.Candle, error)

	// GetCandles returns candles for a given pair and interval
	// Use startDate and endDate to filter candles by date
	// Return a limited count of candles defined by default in the sdk
	// Return hasMore = true if there are more candles to fetch using nextCursor
	// Return nextCursor = the last candle date if there are more candles to fetch
	GetCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time) (res *[]candles.Candle, hasMore bool, nextCursor *time.Time, err error)

	// UpdateCandleListRSI updates only the RSI for a list of candles
	// It returns the updated candles
	UpdateCandleListRSI(ctx context.Context, candles *[]candles.Candle) (*[]candles.Candle, error)

	// QuerySurroundingDates returns the first and last candle date for a given pair and interval
	// It returns 404 not found if no candles are found for the given pair and interval
	QuerySurroundingDates(ctx context.Context, pair common.Pair, interval common.Interval) (*candles.Date, *candles.Date, error)
}

type client struct {
	*sdk.Client
}

func NewCandlesClient(baseURL string) Candles {
	return &client{
		sdk.NewSDKClient(baseURL),
	}
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

	chuncks := createCandlesChunk(newCandles, chunckSize)
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

		createdCandlesChunk := []candles.Candle{}
		err = json.Unmarshal(res, &createdCandlesChunk)
		if err != nil {
			return nil, errors.NewUnexpected("create failed to unmarshal candles while createChunck", err)
		}

		createdCandles = append(createdCandles, createdCandlesChunk...)
	}

	return &createdCandles, nil
}

func (c *client) UpdateCandleListRSI(ctx context.Context, candlesRSIs *[]candles.Candle) (*[]candles.Candle, error) {
	if candlesRSIs == nil || len(*candlesRSIs) == 0 {
		return nil, nil
	}

	updatedCandles := []candles.Candle{}
	chuncks := createCandlesChunk(candlesRSIs, defaultCreateCandlesChunckSize)
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

func (c *client) GetCandles(ctx context.Context, pair common.Pair, interval common.Interval, startDate *time.Time) (*[]candles.Candle, bool, *time.Time, error) {
	queryValues := url.Values{}

	queryValues.Add("pair", string(pair))
	queryValues.Add("interval", string(interval))
	queryValues.Add("limit", strconv.Itoa(defaultGetCandlesLimit))

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
