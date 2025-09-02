package sdk

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	buySignals "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/common"
	"github.com/sopial42/bifrost/pkg/errors"
)

const defaultCreateBuySignalsChunckSize = 1000

type BuySignals interface {
	CreateBuySignals(ctx context.Context, buySignal *[]buySignals.Details) (*[]buySignals.Details, error)
	GetBuySignals(context.Context, common.Pair, common.Interval, buySignals.Name, *time.Time) (res *[]buySignals.Details, hasMore bool, nextCursor *time.Time, err error)
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

func (c *client) GetBuySignals(ctx context.Context, pair common.Pair, interval common.Interval, name buySignals.Name, firstDate *time.Time) (*[]buySignals.Details, bool, *time.Time, error) {
	queryValues := url.Values{}

	queryValues.Add("pair", pair.String())
	queryValues.Add("interval", interval.String())
	queryValues.Add("name", string(name))

	if firstDate != nil {
		queryValues.Add("first_date", firstDate.Format(time.RFC3339))
	}

	res, err := c.Get(ctx, "/buy_signals?"+queryValues.Encode())
	if err != nil {
		return nil, false, nil, err
	}

	getResponse := struct {
		BuySignals []buySignals.Details `json:"buy_signals"`
		HasMore    bool                 `json:"has_more"`
		NextCursor *time.Time           `json:"next_cursor"`
	}{}

	err = json.Unmarshal(res, &getResponse)
	if err != nil {
		return nil, false, nil, errors.NewUnexpected("failed to unmarshal buySignals", err)
	}

	return &getResponse.BuySignals, getResponse.HasMore, getResponse.NextCursor, nil
}
