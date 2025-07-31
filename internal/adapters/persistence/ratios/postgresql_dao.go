package ratios

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	positionsDAO "github.com/sopial42/bifrost/internal/adapters/persistence/positions"
	buySignals "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/common"
	"github.com/sopial42/bifrost/pkg/domains/positions"
	domain "github.com/sopial42/bifrost/pkg/domains/ratios"
)

type RatioDAO struct {
	bun.BaseModel `bun:"table:ratios"`

	ID         uuid.UUID                 `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	SerialID   int64                     `bun:"serial_id,autoincrement"`
	PositionID uuid.UUID                 `bun:"type:uuid"`
	Position   *positionsDAO.PositionDAO `bun:"rel:belongs-to,join:position_id=id"`
	Ratio      float64                   `bun:"ratio"`
	Date       time.Time                 `bun:"date"`
}

func ratiosDetailsToRatiosDAOs(ratios *[]domain.Ratio) []RatioDAO {
	ratiosDAOs := make([]RatioDAO, len(*ratios))
	for i, r := range *ratios {
		ratiosDAOs[i] = RatioDAO{
			ID:         uuid.New(),
			PositionID: uuid.UUID(r.PositionID),
			Ratio:      r.Ratio,
			Date:       r.Date,
		}
	}

	return ratiosDAOs
}

func ratiosDAOsToRatiosDetails(ratiosDAO []RatioDAO) (*[]domain.Ratio, error) {
	ratiosDetails := make([]domain.Ratio, len(ratiosDAO))
	for i, r := range ratiosDAO {
		ratiosDetails[i] = domain.Ratio{
			ID:         r.ID,
			PositionID: positions.ID(r.PositionID),
			Ratio:      r.Ratio,
			Date:       r.Date,
		}
		if r.Position != nil {
			ratiosDetails[i].Position = &positions.Details{
				ID:       positions.ID(r.Position.ID),
				SerialID: r.Position.SerialID,
				Name:     positions.Name(r.Position.Name),
				Fullname: positions.Fullname(r.Position.Fullname),
				TP:       r.Position.TP,
				SL:       r.Position.SL,
				Metadata: r.Position.Metadata,
			}

			if r.Position.BuySignal != nil {
				ratiosDetails[i].Position.BuySignal = &buySignals.Details{
					ID:         buySignals.ID(r.Position.BuySignal.ID),
					BusinessID: buySignals.BusinessID(r.Position.BuySignal.BusinessID),
					Pair:       common.Pair(r.Position.BuySignal.Pair),
					Name:       buySignals.Name(r.Position.BuySignal.Name),
					Fullname:   buySignals.Fullname(r.Position.BuySignal.Fullname),
					Date:       buySignals.Date(r.Position.BuySignal.Date),
					Price:      r.Position.BuySignal.Price,
					Metadata:   r.Position.BuySignal.Metadata,
				}
			}
		}
	}

	return &ratiosDetails, nil
}
