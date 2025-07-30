package ratios

import (
	"context"
	"fmt"

	domain "github.com/bifrost/pkg/domains/ratios"
)

type ratiosService struct {
	persistence Persistence
}

func NewRatiosService(persistence Persistence) Service {
	return &ratiosService{
		persistence: persistence,
	}
}

func (b *ratiosService) CreateRatios(ctx context.Context, ratios *[]domain.Ratio) (*[]domain.Ratio, error) {
	bs, err := b.persistence.InsertRatios(ctx, ratios)
	if err != nil {
		return &[]domain.Ratio{}, fmt.Errorf("unable to create ratios: %w", err)
	}

	return bs, nil
}
