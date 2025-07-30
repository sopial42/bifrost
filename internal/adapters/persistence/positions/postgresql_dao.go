package positions

import (
	bsPersistence "github.com/bifrost/internal/adapters/persistence/buySignals"
	bsDomain "github.com/bifrost/pkg/domains/buySignals"
	positions "github.com/bifrost/pkg/domains/positions"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PositionDAO struct {
	bun.BaseModel `bun:"table:positions"`

	ID          uuid.UUID                   `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	SerialID    int64                       `bun:"serial_id,autoincrement"`
	BuySignalID uuid.UUID                   `bun:"type:uuid"`
	BuySignal   *bsPersistence.BuySignalDAO `bun:"rel:belongs-to,join:buy_signal_id=id"`
	Name        string                      `bun:"name"`
	Fullname    string                      `bun:"fullname"`
	TP          float64                     `bun:"tp"`
	SL          float64                     `bun:"sl"`
	Metadata    map[string]any              `bun:"metadata,type:jsonb"`
}

func positionDetailsToPositionDAOs(positions *[]positions.Details) []PositionDAO {
	positionDAOs := make([]PositionDAO, len(*positions))

	for i, pos := range *positions {
		tp := pos.TP
		sl := pos.SL

		if tp > 0 {
			tp = float64(int(tp*100)) / 100
		}
		if sl > 0 {
			sl = float64(int(sl*100)) / 100
		}

		positionDAOs[i] = PositionDAO{
			ID:          uuid.New(),
			BuySignalID: uuid.UUID(pos.BuySignalID),
			Name:        string(pos.Name),
			Fullname:    string(pos.Fullname),
			TP:          tp,
			SL:          sl,
			Metadata:    pos.Metadata,
		}
	}

	return positionDAOs
}

func positionDAOsToPositionDetails(positionsDAO []PositionDAO) (*[]positions.Details, error) {
	res := make([]positions.Details, len(positionsDAO))
	if len(positionsDAO) == 0 {
		return nil, nil
	}

	for i, p := range positionsDAO {
		res[i] = positions.Details{
			ID:          positions.ID(p.ID),
			SerialID:    p.SerialID,
			BuySignalID: bsDomain.ID(p.BuySignalID),
			Name:        positions.Name(p.Name),
			Fullname:    positions.Fullname(p.Fullname),
			TP:          p.TP,
			SL:          p.SL,
			Metadata:    p.Metadata,
		}

		if p.BuySignal != nil {
			bs := &bsDomain.Details{
				ID:       bsDomain.ID(p.BuySignalID),
				Pair:     p.BuySignal.Pair,
				Date:     bsDomain.Date(p.BuySignal.Date),
				Name:     p.BuySignal.Name,
				Fullname: p.BuySignal.Fullname,
				Price:    p.BuySignal.Price,
				Metadata: p.BuySignal.Metadata,
			}

			res[i].BuySignal = bs
		}
	}

	return &res, nil
}
