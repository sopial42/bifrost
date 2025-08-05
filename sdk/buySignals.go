package sdk

import (
	"context"
	"encoding/json"

	buySignals "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/errors"
)

const defaultCreateBuySignalsChunckSize = 1000

type BuySignals interface {
	CreateBuySignals(ctx context.Context, buySignal *[]buySignals.Details) (*[]buySignals.Details, error)
}

func (c *client) CreateBuySignals(ctx context.Context, newBS *[]buySignals.Details) (*[]buySignals.Details, error) {
	if newBS == nil || len(*newBS) == 0 {
		return nil, nil
	}

	createdBS := []buySignals.Details{}
	// Chunk the candles list into specified size or 1000 if set to 0
	chuncks := createChunk(newBS, defaultCreateBuySignalsChunckSize)
	if chuncks == nil {
		return nil, nil
	}

	for _, chunk := range *chuncks {
		input := map[string]interface{}{
			"buy_signals": chunk,
		}
		body, err := json.Marshal(input)
		if err != nil {
			return nil, errors.NewUnexpected("failed to marshal candles", err)
		}

		res, err := c.Post(ctx, "/buy_signals", body)
		if err != nil {
			return nil, err
		}

		postReponse := struct {
			BuySignals []buySignals.Details `json:"buy_signals"`
		}{}

		err = json.Unmarshal(res, &postReponse)
		if err != nil {
			return nil, errors.NewUnexpected("create failed to unmarshal buySignals while createChunck", err)
		}

		createdBS = append(createdBS, postReponse.BuySignals...)
	}

	return &createdBS, nil
}
