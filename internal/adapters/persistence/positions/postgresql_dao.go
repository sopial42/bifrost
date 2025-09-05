package positions

import (
	"time"

	"github.com/google/uuid"
	bsPersistence "github.com/sopial42/bifrost/internal/adapters/persistence/buySignals"
	bsDomain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	candlesDomain "github.com/sopial42/bifrost/pkg/domains/candles"
	positions "github.com/sopial42/bifrost/pkg/domains/positions"
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
	RatioValue  *float64                    `bun:"ratio_value,nullzero"`
	RatioDate   *time.Time                  `bun:"ratio_date,nullzero"`
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
			BuySignalID: uuid.UUID(pos.BuySignalID),
			Name:        string(pos.Name),
			Fullname:    string(pos.Fullname),
			TP:          tp,
			SL:          sl,
			Metadata:    pos.Metadata,
		}

		if uuid.UUID(pos.ID) != uuid.Nil {
			positionDAOs[i].ID = uuid.UUID(pos.ID)
		} else {
			positionDAOs[i].ID = uuid.New()
		}

		if pos.Ratio != nil {
			ratioDate := time.Time(pos.Ratio.Date)
			positionDAOs[i].RatioValue = &pos.Ratio.Value
			positionDAOs[i].RatioDate = &ratioDate
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
			SerialID:    positions.SerialID(p.SerialID),
			BuySignalID: bsDomain.ID(p.BuySignalID),
			Name:        positions.Name(p.Name),
			Fullname:    positions.Fullname(p.Fullname),
			TP:          p.TP,
			SL:          p.SL,
			Metadata:    p.Metadata,
		}

		if p.RatioValue != nil {
			res[i].Ratio = &positions.Ratio{
				Value: *p.RatioValue,
			}
		}

		if p.RatioDate != nil {
			ratioDate := candlesDomain.Date(*p.RatioDate)
			res[i].Ratio.Date = ratioDate
		}

		if p.BuySignal != nil {
			id := bsDomain.ID(p.BuySignalID)
			bs := &bsDomain.Details{
				ID:       &id,
				Pair:     p.BuySignal.Pair,
				Interval: p.BuySignal.Interval,
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
