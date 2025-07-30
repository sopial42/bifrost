package ratios

import (
	"context"

	domain "github.com/bifrost/pkg/domains/ratios"
)

type Service interface {
	CreateRatios(context.Context, *[]domain.Ratio) (*[]domain.Ratio, error)
}

type Persistence interface {
	InsertRatios(context.Context, *[]domain.Ratio) (*[]domain.Ratio, error)
}
