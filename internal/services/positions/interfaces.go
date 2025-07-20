package positions

import (
	"context"

	domain "github.com/bifrost/internal/domains/positions"
)

type Service interface {
	CreatePositions(context.Context, *[]domain.Details) (*[]domain.Details, error)
}

type Persistence interface {
	InsertPositions(context.Context, *[]domain.Details) (*[]domain.Details, error)
}
