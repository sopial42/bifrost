package positions

import (
	"context"

	domain "github.com/sopial42/bifrost/pkg/domains/positions"
)

type positionsService struct {
	persistence Persistence
}

func NewPositionsService(persistence Persistence) Service {
	return &positionsService{
		persistence: persistence,
	}
}

func (p *positionsService) CreatePositions(ctx context.Context, positions *[]domain.Details) (*[]domain.Details, error) {
	pos, err := p.persistence.InsertPositions(ctx, positions)
	if err != nil {
		return &[]domain.Details{}, err
	}

	return pos, nil
}
