package positions

import (
	"context"

	domain "github.com/sopial42/bifrost/pkg/domains/positions"
)

type Service interface {
	CreatePositions(context.Context, *[]domain.Details) (*[]domain.Details, error)
	ComputeAllPositions(context.Context) (int, error)
}

type Persistence interface {
	InsertPositions(context.Context, *[]domain.Details) (*[]domain.Details, error)
	InsertRatios(context.Context, *[]domain.Details) (*[]domain.Details, error)
	GetPositionsWithNoRatio(ctx context.Context, cursor *int64, limit int) (positions *[]domain.Details, hasMore bool, nextCursor *int64, err error)
}
