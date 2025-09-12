package buysignals

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"

	buySignalsSVC "github.com/sopial42/bifrost/internal/services/buySignals"
	domain "github.com/sopial42/bifrost/pkg/domains/buySignals"
	"github.com/sopial42/bifrost/pkg/domains/common"
	appErrors "github.com/sopial42/bifrost/pkg/errors"
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
		return &[]domain.Details{}, fmt.Errorf("unable to insert buy signal, nil or empty")
	}

	buySignalsDAO := buySignalDetailsToBuySignalDAOs(ctx, bsReports, false)
	_, err := c.clientDB.
		NewInsert().
		Model(&buySignalsDAO).
		Exec(ctx)
	if err != nil {
		if errPg, ok := err.(pgdriver.Error); ok && errPg.Field('C') == pgerrcode.UniqueViolation {
			return &[]domain.Details{}, appErrors.NewAlreadyExists("buy signal already exists, unique constraint violation")
		}

		return &[]domain.Details{}, fmt.Errorf("unable to insert buy signals: %w", err)
	}

	res := buySignalDAOsToBuySignalDetails(ctx, &buySignalsDAO)
	log.Debugf("Insert buySignals done")
	return res, nil
}

func (c *pgPersistence) QueryBuySignals(ctx context.Context, pair common.Pair, interval common.Interval, name domain.Name, firstDate *time.Time, limit int) (*[]domain.Details, bool, *time.Time, error) {
	buySignalsDAO := []BuySignalDAO{}
	request := c.clientDB.NewSelect().
		Model(&buySignalsDAO).
		Where("pair = ?", pair).
		Where("interval = ?", interval).
		Where("name = ?", name).
		OrderExpr("date ASC")

	if limit > 0 {
		request.Limit(limit + 1)
	}

	if firstDate != nil && !firstDate.IsZero() {
		request.Where("date >= ?", firstDate)
	}

	if limit > 0 {
		request.Limit(limit + 1)
	}

	err := request.Scan(ctx)
	if err != nil {
		return nil, false, nil, fmt.Errorf("unable to perform db query: %v", err)
	}

	if limit <= 0 {
		return buySignalDAOsToBuySignalDetails(ctx, &buySignalsDAO), false, nil, nil
	}

	hasMore := len(buySignalsDAO) > limit
	var nextCursor *time.Time
	if hasMore {
		nextCursor = &buySignalsDAO[len(buySignalsDAO)-1].Date
		buySignalsDAO = buySignalsDAO[:limit]
	}

	return buySignalDAOsToBuySignalDetails(ctx, &buySignalsDAO), hasMore, nextCursor, nil
}

func (c *pgPersistence) UpsertBuySignals(ctx context.Context, bs domain.Details) (*[]domain.Details, error) {
	bsDAO := buySignalDetailsToBuySignalDAOs(ctx, &[]domain.Details{bs}, true)
	if len(bsDAO) == 0 {
		return &[]domain.Details{}, fmt.Errorf("unable to upsert buy signals, empty")
	}

	_, err := c.clientDB.
		NewInsert().
		Model(&bsDAO).
		On("CONFLICT (pair, interval, fullname, business_id) DO UPDATE").
		Set("price = EXCLUDED.price").
		Set("metadata = EXCLUDED.metadata").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return &[]domain.Details{}, fmt.Errorf("unable to upsert buy signals: %w", err)
	}

	return buySignalDAOsToBuySignalDetails(ctx, &bsDAO), nil
}
