package storage

import (
	"context"
	"fmt"

	"github.com/NStegura/saga/golibs/event"
	eventRepo "github.com/NStegura/saga/golibs/event/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	orderRepo "github.com/NStegura/saga/orders/internal/storage/repo/order"
	stateRepo "github.com/NStegura/saga/orders/internal/storage/repo/state"
)

type DB struct {
	orderRepo.OrderRepo
	eventRepo.EventRepo
	stateRepo.StateRepo
	pool *pgxpool.Pool

	logger *logrus.Logger
}

func New(
	ctx context.Context,
	dsn string,
	logger *logrus.Logger,
	runMigrations bool,
) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}

	db := DB{
		OrderRepo: orderRepo.New(logger, pool),
		EventRepo: event.NewEventRepository(logger),
		StateRepo: stateRepo.New(logger, pool),
		pool:      pool,
		logger:    logger,
	}
	if !runMigrations {
		return &db, nil
	}

	if err = db.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	return &db, nil
}

func (db *DB) Shutdown(_ context.Context) {
	db.logger.Debug("db shutdown")
	db.pool.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	db.logger.Debug("Ping db")
	err := db.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("DB ping eror, %w", err)
	}
	return nil
}
