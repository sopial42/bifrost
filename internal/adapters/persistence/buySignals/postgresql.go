package buysignals

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"

	buySignalsSVC "github.com/sopial42/bifrost/internal/services/buySignals"
	domain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/logger"
)

type pgPersistence struct {
	clientDB *bun.DB
}

func NewPersistence(client *bun.DB) buySignalsSVC.Persistence {
	return &pgPersistence{clientDB: client}
}

func (c *pgPersistence) InsertBuySignals(ctx context.Context, bsReports *[]domain.Details) (*[]domain.Details, error) {
	log := logger.GetLogger(ctx)
	if bsReports == nil || len(*bsReports) == 0 {
		return &[]domain.Details{}, nil
	}

	buySignalsDAO := buySignalDetailsToBuySignalDAOs(bsReports)
	_, err := c.clientDB.
		NewInsert().
		Model(&buySignalsDAO).
		Ignore().
		Exec(ctx)
	if err != nil {
		if errPg, ok := err.(pgdriver.Error); ok && errPg.Field('C') == pgerrcode.UniqueViolation {
			return &[]domain.Details{}, nil
		}

		log.Errorf("Query refused, throw pgError, %v", err)
		return &[]domain.Details{}, err
	}

	res := buySignalDAOsToBuySignalDetails(buySignalsDAO)
	log.Debugf("Insert buySignals done")
	return res, nil
}
