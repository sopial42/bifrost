package positions

import (
	"context"

	domain "github.com/sopial42/bifrost/pkg/domains/positions"
)

type Service interface {
	CreatePositions(context.Context, *[]domain.Details) (*[]domain.Details, error)
}

type Persistence interface {
	InsertPositions(context.Context, *[]domain.Details) (*[]domain.Details, error)
}
