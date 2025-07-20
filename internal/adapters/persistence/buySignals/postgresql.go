package buysignals

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	domain "github.com/bifrost/internal/domains/buySignals"
	"github.com/bifrost/internal/domains/common"
)

type BuySignalDAO struct {
	bun.BaseModel `bun:"table:buy_signals"`

	ID         uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	BusinessID domain.BusinessID
	Name       domain.Name
	Fullname   domain.Fullname
	Pair       common.Pair
	Interval   common.Interval
	Date       time.Time
	Price      float64
	Metadata   map[string]any `bun:"metadata,type:jsonb,nullzero"`
}
