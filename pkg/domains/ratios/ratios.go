package ratios

import (
	"time"

	"github.com/bifrost/pkg/domains/positions"
	"github.com/google/uuid"
)

type Ratio struct {
	ID         uuid.UUID          `json:"id"`
	PositionID positions.ID       `json:"position_id"`
	Position   *positions.Details `json:"position,omitempty"`
	Ratio      float64            `json:"ratio"`
	Date       time.Time          `json:"date"`
}
