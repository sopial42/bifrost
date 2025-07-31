package ratios

import (
	"time"

	"github.com/google/uuid"
	"github.com/sopial42/bifrost/pkg/domains/positions"
)

type Ratio struct {
	ID         uuid.UUID          `json:"id"`
	PositionID positions.ID       `json:"position_id"`
	Position   *positions.Details `json:"position,omitempty"`
	Ratio      float64            `json:"ratio"`
	Date       time.Time          `json:"date"`
}
