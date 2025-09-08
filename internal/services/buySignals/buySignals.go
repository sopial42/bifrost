package buySignals

import (
	"context"
	"fmt"
	"time"

	domain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/common"
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

func (b *buySignalsService) GetBuySignals(ctx context.Context, pair common.Pair, interval common.Interval, name domain.Name, date *time.Time, limit int) (*[]domain.Details, bool, *time.Time, error) {
	bs, hasMore, nextCursor, err := b.persistence.QueryBuySignals(ctx, pair, interval, name, date, limit)
	if err != nil {
		return &[]domain.Details{}, false, nil, fmt.Errorf("unable to get buy signals: %w", err)
	}

	return bs, hasMore, nextCursor, nil
}

func (b *buySignalsService) UpsertBuySignals(ctx context.Context, buySignal domain.Details) (*[]domain.Details, error) {
	bs, err := b.persistence.UpsertBuySignals(ctx, buySignal)
	if err != nil {
		return &[]domain.Details{}, fmt.Errorf("unable to upsert buy signals: %w", err)
	}

	return bs, nil
}
