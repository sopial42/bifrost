package buySignals

import (
	"context"
	"fmt"

	domain "github.com/bifrost/pkg/domains/buySignals"
)

type buySignalsService struct {
	persistence Persistence
}

func NewBuySignalsService(persistence Persistence) Service {
	return &buySignalsService{
		persistence: persistence,
	}
}

func (b *buySignalsService) CreateBuySignals(ctx context.Context, buySignals *[]domain.Details) (*[]domain.Details, error) {
	bs, err := b.persistence.InsertBuySignals(ctx, buySignals)
	if err != nil {
		return &[]domain.Details{}, fmt.Errorf("unable to create buy signals: %w", err)
	}

	return bs, nil
}
