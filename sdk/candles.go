package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/bifrost/pkg/domains/candles"
	"github.com/bifrost/pkg/domains/common"
	"github.com/bifrost/pkg/errors"
	"github.com/bifrost/pkg/sdk"
)

const defaultChunckSize = 500

type Candles interface {
	// CreateCandles insert candles in the database, if a candle already exists, it will be ignored
	// It returns only the newly inserted candles
	// The candle list is chunked by specified size or defaultChunckSize if set to <= 0
	CreateCandles(candles *[]candles.Candle, chunckSize int) (*[]candles.Candle, error)

	// QuerySurroundingDates returns the first and last candle date for a given pair and interval
	// It returns 404 not found if no candles are found for the given pair and interval
	QuerySurroundingDates(pair common.Pair, interval common.Interval) (*candles.Date, *candles.Date, error)
}

type client struct {
	*sdk.Client
}

func NewCandlesClient(baseURL string) Candles {
	fmt.Printf("NewCandlesClient %s\n", baseURL)
	fmt.Printf("NewCandlesClient %s\n", baseURL)
	fmt.Printf("NewCandlesClient %s\n", baseURL)
	fmt.Printf("NewCandlesClient %s\n", baseURL)
	return &client{
		sdk.NewSDKClient(baseURL),
	}
}

func (c *client) CreateCandles(newCandles *[]candles.Candle, chunckSize int) (*[]candles.Candle, error) {
	if len(*newCandles) == 0 {
		return nil, nil
	}

	createdCandles := []candles.Candle{}
	// Chunk the candles list into specified size or 1000 if set to 0
	if chunckSize <= 0 {
		chunckSize = defaultChunckSize
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

		res, err := c.Post("/candles", body)
		if err != nil {
			return nil, err
		}

		createdCandlesChunk := []candles.Candle{}
		err = json.Unmarshal(res, &createdCandlesChunk)
		if err != nil {
			return nil, errors.NewUnexpected("failed to unmarshal candles: "+string(res), err)
		}

		createdCandles = append(createdCandles, createdCandlesChunk...)
	}

	return &createdCandles, nil
}

func (c *client) QuerySurroundingDates(pair common.Pair, interval common.Interval) (*candles.Date, *candles.Date, error) {
	res, err := c.Get("/candles/surrounding-dates?pair=" + string(pair) + "&interval=" + string(interval))
	if err != nil {
		return nil, nil, err
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
