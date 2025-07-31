package buysignals

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	domain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/common"
)

type BuySignalDAO struct {
	bun.BaseModel `bun:"table:buy_signals"`

	ID         uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	BusinessID domain.BusinessID
	Name       domain.Name
	Fullname   domain.Fullname
	Pair       common.Pair
	Date       time.Time
	Price      float64
	Metadata   map[string]any `bun:"metadata,type:jsonb,nullzero"`
}

func buySignalDetailsToBuySignalDAOs(buySignals *[]domain.Details) []BuySignalDAO {
	buySignalsDAO := make([]BuySignalDAO, len(*buySignals))

	for i, bs := range *buySignals {
		buySignalsDAO[i] = BuySignalDAO{
			ID:         uuid.UUID(bs.ID),
			BusinessID: domain.BusinessID(bs.BusinessID),
			Name:       domain.Name(bs.Name),
			Fullname:   domain.Fullname(bs.Fullname),
			Pair:       common.Pair(bs.Pair),
			Date:       time.Time(bs.Date),
			Price:      bs.Price,
			Metadata:   bs.Metadata,
		}
	}

	return buySignalsDAO
}

func buySignalDAOsToBuySignalDetails(buySignalsDAO []BuySignalDAO) *[]domain.Details {
	buySignals := make([]domain.Details, len(buySignalsDAO))

	for i, bs := range buySignalsDAO {
		buySignals[i] = domain.Details{
			ID:         domain.ID(bs.ID),
			BusinessID: domain.BusinessID(bs.BusinessID),
			Name:       domain.Name(bs.Name),
			Fullname:   domain.Fullname(bs.Fullname),
			Pair:       common.Pair(bs.Pair),
			Date:       domain.Date(bs.Date),
			Price:      bs.Price,
			Metadata:   bs.Metadata,
		}
	}

	return &buySignals
}
