package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/sopial42/bifrost/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type PGClient struct {
	Client *bun.DB
}

func NewPGClient(cfg config.DBConfig) *PGClient {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
		pgdriver.WithUser(cfg.User),
		pgdriver.WithPassword(cfg.Password),
		pgdriver.WithDatabase(cfg.DBName),
		pgdriver.WithInsecure(cfg.Unsecure),
		pgdriver.WithTimeout(5*time.Second),
	))

	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)
	sqldb.SetConnMaxLifetime(30 * time.Minute)

	err := sqldb.Ping()
	if err != nil {
		panic(err)
	}
	client := bun.NewDB(sqldb, pgdialect.New())
	client.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("DB_LOG_LEVEL"),
	))

	// client.AddQueryHook(tracing.NewTracingHook(cfg.TracingEnabled))
	return &PGClient{Client: client}
}

func (pg *PGClient) Ping(c context.Context) error {
	err := pg.Client.PingContext(c)
	if err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}
	return nil
}
