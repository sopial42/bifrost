package sdk

import (
	"context"
	"encoding/json"

	"github.com/sopial42/bifrost/pkg/domains/positions"
	"github.com/sopial42/bifrost/pkg/errors"
	"github.com/sopial42/bifrost/pkg/logger"
)

const defaultCreatePositionsChunckSize = 1

type Positions interface {
	CreatePositions(ctx context.Context, positions *[]positions.Details, chunckSize int) (*[]positions.Details, error)
}

func (c *client) CreatePositions(ctx context.Context, newPositions *[]positions.Details, chunckSize int) (*[]positions.Details, error) {
	if newPositions == nil || len(*newPositions) == 0 {
		return nil, nil
	}

	createdPositions := []positions.Details{}
	// Chunk the candles list into specified size or 1000 if set to 0
	if chunckSize <= 0 {
		chunckSize = defaultCreatePositionsChunckSize
	}

	log := logger.GetLogger(ctx)
	log.Infof("Creating %d positions", len(*newPositions))
	chuncks := createChunk(newPositions, chunckSize)
	if chuncks == nil {
		return nil, nil
	}

	for _, chunk := range *chuncks {
		input := map[string]interface{}{
			"positions": chunk,
		}

		body, err := json.Marshal(input)
		if err != nil {
			return nil, errors.NewUnexpected("failed to marshal positions", err)
		}

		res, err := c.Post(ctx, "/positions", body)
		if err != nil {
			return nil, err
		}

		postResponse := struct {
			Positions []positions.Details `json:"positions"`
		}{}

		err = json.Unmarshal(res, &postResponse)
		if err != nil {
			return nil, errors.NewUnexpected("create failed to unmarshal positions while createChunck", err)
		}

		createdPositions = append(createdPositions, postResponse.Positions...)
	}

	log.Infof("Created %d positions", len(createdPositions))
	return &createdPositions, nil
}
