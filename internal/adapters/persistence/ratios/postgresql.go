package ratios

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"

	ratiosSVC "github.com/sopial42/bifrost/internal/services/ratios"
	domain "github.com/sopial42/bifrost/pkg/domains/ratios"
	"github.com/sopial42/bifrost/pkg/logger"
)

type pgPersistence struct {
	clientDB *bun.DB
}

func NewPersistence(client *bun.DB) ratiosSVC.Persistence {
	return &pgPersistence{clientDB: client}
}

func (p *pgPersistence) InsertRatios(ctx context.Context, ratios *[]domain.Ratio) (*[]domain.Ratio, error) {
	log := logger.GetLogger(ctx)
	if ratios == nil || len(*ratios) == 0 {
		log.Warnf("unable to insert, ratios details is nil, or len = 0: %v", ratios)
		return &[]domain.Ratio{}, nil
	}

	ratioDAOs := ratiosDetailsToRatiosDAOs(ratios)
	_, err := p.clientDB.
		NewInsert().
		Model(&ratioDAOs).
		Ignore().
		Returning("*").
		Exec(ctx)
	if err != nil {
		if errPg, ok := err.(pgdriver.Error); ok && errPg.Field('C') == pgerrcode.UniqueViolation {
			return &[]domain.Ratio{}, nil
		}

		log.Errorf("Query refused, throw pgError, %v", err)
		return &[]domain.Ratio{}, err
	}

	log.Debugf("Insert %d ratios done", len(ratioDAOs))
	return ratiosDAOsToRatiosDetails(ratioDAOs)
}
