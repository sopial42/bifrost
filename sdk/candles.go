package sdk

import (
	"encoding/json"

	"github.com/bifrost/pkg/domains/candles"
	"github.com/bifrost/pkg/errors"
	"github.com/bifrost/pkg/sdk"
)

type Candles interface {
	// CreateCandles insert candles in the database, if a candle already exists, it will be ignored
	// It returns only the newly inserted candles
	CreateCandles(candles *[]candles.Candle) (*[]candles.Candle, error)
}

type client struct {
	*sdk.Client
}

func NewCandlesClient(baseURL string) Candles {
	return &client{
		sdk.NewSDKClient(baseURL),
	}
}

func (c *client) CreateCandles(newCandles *[]candles.Candle) (*[]candles.Candle, error) {
	input := map[string]interface{}{
		"candles": newCandles,
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, errors.NewUnexpected("failed to marshal candles", err)
	}

	res, err := c.Post("/candles", body)
	if err != nil {
		return nil, err
	}

	createdCandles := []candles.Candle{}
	err = json.Unmarshal(res, &createdCandles)
	if err != nil {
		return nil, errors.NewUnexpected("failed to unmarshal candles", err)
	}

	return &createdCandles, nil
}
