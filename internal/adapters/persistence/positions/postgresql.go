package positions

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"

	domain "github.com/bifrost/internal/domains/positions"
	"github.com/bifrost/pkg/logger"
	positionSVC "github.com/bifrost/internal/services/positions"
)

type pgPersistence struct {
	clientDB *bun.DB
}

func NewPersistence(client *bun.DB) positionSVC.Persistence {
	return &pgPersistence{clientDB: client}
}

func (p *pgPersistence) InsertPositions(ctx context.Context, pos *[]domain.Details) (*[]domain.Details, error) {
	log := logger.GetLogger(ctx)
	if pos == nil || len(*pos) == 0 {
		log.Warnf("unable to insert, position details is nil, or len = 0: %v", pos)
		return &[]domain.Details{}, nil
	}

	positionDAOs := positionDetailsToPositionDAOs(pos)
	_, err := p.clientDB.
		NewInsert().
		Model(&positionDAOs).
		Ignore().
		Exec(ctx)
	if err != nil {
		if errPg, ok := err.(pgdriver.Error); ok && errPg.Field('C') == pgerrcode.UniqueViolation {
			return &[]domain.Details{}, nil
		}

		log.Errorf("Query refused, throw pgError, %v", err)
		return &[]domain.Details{}, err
	}

	log.Debugf("Insert positions done")
	return positionDAOsToPositionDetails(positionDAOs)
}
