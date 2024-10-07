package repo

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/golibs/event/repo"
	"github.com/NStegura/saga/payments/internal/repo/payment"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type DB struct {
	payment.PaymentRepo
	repo.EventRepo
	pool *pgxpool.Pool

	logger *logrus.Logger
}

func New(
	ctx context.Context,
	dsn string,
	paymentRepo payment.PaymentRepo,
	eventRepo repo.EventRepo,
	logger *logrus.Logger,
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
		PaymentRepo: paymentRepo,
		EventRepo:   eventRepo,
		pool:        pool,
		logger:      logger,
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
