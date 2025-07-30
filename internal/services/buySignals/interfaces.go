package buySignals

import (
	"context"

	domain "github.com/bifrost/pkg/domains/buySignals"
)

type Service interface {
	CreateBuySignals(context.Context, *[]domain.Details) (*[]domain.Details, error)
}

type Persistence interface {
	InsertBuySignals(context.Context, *[]domain.Details) (*[]domain.Details, error)
}
