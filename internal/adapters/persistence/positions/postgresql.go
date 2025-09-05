package positions

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"

	positionSVC "github.com/sopial42/bifrost/internal/services/positions"
	domain "github.com/sopial42/bifrost/pkg/domains/positions"
	"github.com/sopial42/bifrost/pkg/logger"
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

func (p *pgPersistence) InsertRatios(ctx context.Context, pos *[]domain.Details) (*[]domain.Details, error) {
	log := logger.GetLogger(ctx)
	if pos == nil || len(*pos) == 0 {
		log.Debugf("unable to insert, position details is nil, or len = 0: %v", pos)
		return &[]domain.Details{}, nil
	}

	positionDAOs := positionDetailsToPositionDAOs(pos)
	var res []PositionDAO
	err := p.clientDB.
		NewUpdate().
		Model(&positionDAOs).
		Column("ratio_value", "ratio_date").
		Bulk().
		Returning("position_dao.*").
		Scan(ctx, &res)
	if err != nil {
		return &[]domain.Details{}, err
	}

	return positionDAOsToPositionDetails(res)
}

func (p *pgPersistence) GetPositionsWithNoRatio(ctx context.Context, cursor *int64, limit int) (positions *[]domain.Details, hasMore bool, nextCursor *int64, err error) {
	positionsDAO := []PositionDAO{}
	request := p.clientDB.NewSelect().Model(&positionsDAO).
		Where("ratio_value IS NULL").
		Relation("BuySignal").
		OrderExpr("serial_id ASC")

	if cursor != nil {
		request.Where("serial_id >= ?", *cursor)
	}

	if limit > 0 {
		request.Limit(limit + 1)
	}

	err = request.Scan(ctx)
	if err != nil {
		return nil, false, nil, fmt.Errorf("unable to perform db query: %v", err)
	}

	if limit <= 0 {
		positionsModel, err := positionDAOsToPositionDetails(positionsDAO)
		if err != nil {
			return nil, false, nil, fmt.Errorf("unable to convert positionsDAO to positionsModel and no limit: %w", err)
		}

		return positionsModel, false, nil, nil
	}

	hasMore = len(positionsDAO) > limit
	if hasMore {
		nextCursor = &positionsDAO[len(positionsDAO)-1].SerialID
		positionsDAO = positionsDAO[:limit]
	}

	positionsModel, err := positionDAOsToPositionDetails(positionsDAO)
	if err != nil {
		return nil, false, nil, fmt.Errorf("unable to convert positionsDAO to positionsModel: %w", err)
	}

	return positionsModel, hasMore, nextCursor, err
}

func (p *pgPersistence) GetPositionByID(ctx context.Context, id domain.ID) (*domain.Details, error) {
	positionDAO := PositionDAO{}
	err := p.clientDB.NewSelect().Model(&positionDAO).
		Relation("BuySignal").
		Where("position_dao.id = ?", id.String()).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	positionModel, err := positionDAOsToPositionDetails([]PositionDAO{positionDAO})
	if err != nil {
		return nil, err
	}

	if positionModel == nil || len(*positionModel) == 0 {
		return nil, nil
	}

	return &(*positionModel)[0], nil
}
