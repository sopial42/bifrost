package inProcess

import (
	persistence "github.com/sopial42/bifrost/pkg/adapters/persistence"
	buySignalsPersistence "github.com/sopial42/bifrost/pkg/adapters/persistence/buySignals"
	candlesPersistence "github.com/sopial42/bifrost/pkg/adapters/persistence/candles"
	positionsPersistence "github.com/sopial42/bifrost/pkg/adapters/persistence/positions"
	"github.com/sopial42/bifrost/pkg/common/config"
	"github.com/sopial42/bifrost/pkg/ports"
	buySignalsSVC "github.com/sopial42/bifrost/pkg/services/buySignals"
	candlesSVC "github.com/sopial42/bifrost/pkg/services/candles"
	positionsSVC "github.com/sopial42/bifrost/pkg/services/positions"
)

type inProcessClient struct {
	buySignalsSVC buySignalsSVC.Service
	candlesSVC    candlesSVC.Service
	positionsSVC  positionsSVC.Service
}

func NewBifrostInProcessClient(dbConfig config.DBConfig) ports.Client {
	pgClient := persistence.NewPGClient(dbConfig)
	buySignalsPersistence := buySignalsPersistence.NewPersistence(pgClient.Client)
	candlesPersistence := candlesPersistence.NewPersistence(pgClient.Client)
	positionsPersistence := positionsPersistence.NewPersistence(pgClient.Client)

	buySignalsSVC := buySignalsSVC.NewBuySignalsService(buySignalsPersistence)
	candlesSVC := candlesSVC.NewCandlesService(candlesPersistence)
	positionsSVC := positionsSVC.NewPositionsService(positionsPersistence, candlesSVC, buySignalsSVC)

	return &inProcessClient{
		buySignalsSVC: buySignalsSVC,
		candlesSVC:    candlesSVC,
		positionsSVC:  positionsSVC,
	}
}
