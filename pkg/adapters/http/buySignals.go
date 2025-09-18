package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sopial42/bifrost/pkg/common/logger"
	"github.com/sopial42/bifrost/pkg/common/sdk"
	domain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/common"
	appErrors "github.com/sopial42/bifrost/pkg/common/errors"
)

const defaultCreateBuySignalsChunckSize = 1000
const defaultGetBuySignalsLimit = 1000



func (c *client) CreateBuySignals(ctx context.Context, newBS *[]domain.Details) (*[]domain.Details, error) {
	log := logger.GetLogger(ctx)

	if newBS == nil || len(*newBS) == 0 {
		return nil, nil
	}

	log.Infof("Creating %d buy signals", len(*newBS))
	createdBS := []domain.Details{}
	chuncks := sdk.CreateChunk(newBS, defaultCreateBuySignalsChunckSize)
	if chuncks == nil {
		return nil, nil
	}

	for _, chunk := range *chuncks {
		input := map[string]interface{}{
			"buy_signals": chunk,
		}
		body, err := json.Marshal(input)
		if err != nil {
			return nil, appErrors.NewUnexpected("failed to marshal candles", err)
		}

		res, err := c.Post(ctx, "/buy_signals", body)
		if err != nil {
			if errors.Is(err, appErrors.ErrAlreadyExists) {
				log.Debugf("buySignals already exists: %+v", err)
			} else {
				return nil, fmt.Errorf("failed to post buySignals: %w", err)
			}
		}

		postReponse := struct {
			BuySignals []domain.Details `json:"buy_signals"`
		}{}

		err = json.Unmarshal(res, &postReponse)
		if err != nil {
			return nil, appErrors.NewUnexpected("create failed to unmarshal buySignals while createChunck", err)
		}

		createdBS = append(createdBS, postReponse.BuySignals...)
	}

	log.Infof("created %d buy signals", len(createdBS))
	return &createdBS, nil
}

func (c *client) GetBuySignals(ctx context.Context, pair common.Pair, interval common.Interval, name domain.Name, firstDate *time.Time) (*[]domain.Details, bool, *time.Time, error) {
	queryValues := url.Values{}

	queryValues.Add("pair", pair.String())
	queryValues.Add("interval", interval.String())
	queryValues.Add("name", string(name))
	queryValues.Add("limit", strconv.Itoa(defaultGetBuySignalsLimit))

	if firstDate != nil {
		queryValues.Add("first_date", firstDate.Format(time.RFC3339))
	}

	res, err := c.Get(ctx, "/buy_signals?"+queryValues.Encode())
	if err != nil {
		return nil, false, nil, err
	}

	getResponse := struct {
		BuySignals []domain.Details `json:"buy_signals"`
		HasMore    bool             `json:"has_more"`
		NextCursor *time.Time       `json:"next_cursor"`
	}{}

	err = json.Unmarshal(res, &getResponse)
	if err != nil {
		return nil, false, nil, appErrors.NewUnexpected("failed to unmarshal buySignals", err)
	}

	return &getResponse.BuySignals, getResponse.HasMore, getResponse.NextCursor, nil
}
