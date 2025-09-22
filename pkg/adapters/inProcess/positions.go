package inProcess

import (
	"context"

	"github.com/sopial42/bifrost/pkg/domains/positions"
)

func (c *inProcessClient) CreatePositions(ctx context.Context, newPositions *[]positions.Details, chunckSize int) (*[]positions.Details, error) {
	return nil, nil
}
