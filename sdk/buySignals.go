package sdk

import (
	"context"

	buySignals "github.com/sopial42/bifrost/pkg/domains/buySignals"
)

type BuySignals interface {
	CreateBuySignals(ctx context.Context, buySignal *[]buySignals.Details) (*[]buySignals.Details, error)
}

func (c *client) CreateBuySignals(ctx context.Context, buySignal *[]buySignals.Details) (*[]buySignals.Details, error) {
	return nil, nil
}
