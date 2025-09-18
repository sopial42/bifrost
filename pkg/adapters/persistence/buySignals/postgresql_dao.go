package buysignals

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/sopial42/bifrost/pkg/common/logger"
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
	Interval   common.Interval
	Date       time.Time
	Price      float64
	Metadata   map[string]any `bun:"metadata,type:jsonb,nullzero"`
}

func buySignalDetailsToBuySignalDAOs(ctx context.Context, buySignals *[]domain.Details, isUpdate bool) []BuySignalDAO {
	log := logger.GetLogger(ctx)

	buySignalsDAO := make([]BuySignalDAO, len(*buySignals))

	for i, bs := range *buySignals {
		buySignalsDAO[i] = BuySignalDAO{
			BusinessID: domain.BusinessID(bs.BusinessID),
			Name:       domain.Name(bs.Name),
			Fullname:   domain.Fullname(bs.Fullname),
			Pair:       common.Pair(bs.Pair),
			Interval:   common.Interval(bs.Interval),
			Date:       time.Time(bs.Date),
			Price:      bs.Price,
			Metadata:   bs.Metadata,
		}

		if isUpdate {
			if bs.ID == nil {
				log.Debugf("unable to update buy signal as ID is nil: %+v", bs)
				continue
			}

			buySignalsDAO[i].ID = uuid.UUID(*bs.ID)
		} else {
			buySignalsDAO[i].ID = uuid.New()
		}
	}

	return buySignalsDAO
}

func buySignalDAOsToBuySignalDetails(ctx context.Context, buySignalsDAO *[]BuySignalDAO) *[]domain.Details {
	log := logger.GetLogger(ctx)
	if buySignalsDAO == nil {
		return &[]domain.Details{}
	}

	buySignals := make([]domain.Details, len(*buySignalsDAO))

	for i, bs := range *buySignalsDAO {
		if bs.ID == uuid.Nil {
			log.Warnf("unable to convert buy signal as ID is nil: %+v", bs)
			continue
		}

		id := domain.ID(bs.ID)
		buySignals[i] = domain.Details{
			ID:         &id,
			BusinessID: domain.BusinessID(bs.BusinessID),
			Name:       domain.Name(bs.Name),
			Fullname:   domain.Fullname(bs.Fullname),
			Pair:       common.Pair(bs.Pair),
			Interval:   common.Interval(bs.Interval),
			Date:       domain.Date(bs.Date),
			Price:      bs.Price,
			Metadata:   bs.Metadata,
		}
	}

	return &buySignals
}
