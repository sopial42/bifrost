package inProcess

import (
	"context"
	"time"

	domain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/common"
)

func (c *inProcessClient) CreateBuySignals(ctx context.Context, newBS *[]domain.Details) (*[]domain.Details, error) {
	return nil, nil
}

func (c *inProcessClient) GetBuySignals(ctx context.Context, pair common.Pair, interval common.Interval, name domain.Name, firstDate *time.Time) (*[]domain.Details, bool, *time.Time, error) {
	return nil, false, nil, nil
}
