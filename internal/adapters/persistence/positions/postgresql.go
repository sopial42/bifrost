package positions

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"

	positionSVC "github.com/sopial42/bifrost/internal/services/positions"
	domain "github.com/sopial42/bifrost/pkg/domains/positions"
	appErrors "github.com/sopial42/bifrost/pkg/errors"
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
		Exec(ctx)
	if err != nil {
		if errPg, ok := err.(pgdriver.Error); ok && errPg.Field('C') == pgerrcode.UniqueViolation {
			return &[]domain.Details{}, appErrors.NewAlreadyExists("position already exists")
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
		Where("tp > 0").
		Where("sl > 0").
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

func (p *pgPersistence) GetPositionsWithNoRatioCount(ctx context.Context) (count int, err error) {
	// Add check on tp and sl != 0
	count, err = p.clientDB.NewSelect().
		Model(&PositionDAO{}).
		Where("ratio_value IS NULL").
		Where("tp > 0").
		Where("sl > 0").
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to perform ratio count db query: %v", err)
	}
	return count, nil
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

func (p *pgPersistence) UpsertPosition(ctx context.Context, position *domain.Details) (*domain.Details, error) {
	log := logger.GetLogger(ctx)
	log.Debugf("Upsert position start with position: %+v", position)

	positionDAO := positionDetailsToPositionDAOs(&[]domain.Details{*position})

	log.Infof("Upsert position DAOs: %+v", positionDAO)
	_, err := p.clientDB.
		NewInsert().
		Model(&positionDAO).
		On("CONFLICT (buy_signal_id, fullname) DO UPDATE").
		Set("tp = EXCLUDED.tp").
		Set("sl = EXCLUDED.sl").
		Column("name", "fullname", "buy_signal_id", "tp", "sl").
		Returning("*").
		Exec(ctx)
	if err != nil || len(positionDAO) == 0 {
		return nil, fmt.Errorf("unable to upsert position: %w", err)
	}

	log.Debugf("positionDAO length: %+v", positionDAO)
	positionModel, err := positionDAOsToPositionDetails([]PositionDAO{positionDAO[0]})
	if err != nil {
		return nil, fmt.Errorf("unable to convert positionDAO to positionModel: %w", err)
	}

	if positionModel == nil || len(*positionModel) == 0 {
		return nil, nil
	}

	return &(*positionModel)[0], nil
}
