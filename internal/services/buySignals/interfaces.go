package buySignals

import (
	"context"
	"time"

	domain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/common"
)

type Service interface {
	CreateBuySignals(context.Context, *[]domain.Details) (*[]domain.Details, error)
	GetBuySignals(context.Context, common.Pair, common.Interval, domain.Name, *time.Time, int) (*[]domain.Details, bool, *time.Time, error)
}

type Persistence interface {
	InsertBuySignals(context.Context, *[]domain.Details) (*[]domain.Details, error)
	QueryBuySignals(context.Context, common.Pair, common.Interval, domain.Name, *time.Time, int) (*[]domain.Details, bool, *time.Time, error)
}
