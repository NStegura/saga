package main

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/golibs/event/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type DB struct {
	repo.EventRepo
	pool *pgxpool.Pool

	logger *logrus.Logger
}

func New(ctx context.Context, dsn string, eventRepo repo.EventRepo, logger *logrus.Logger) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}
	return &DB{
		EventRepo: eventRepo,
		pool:      pool,
		logger:    logger,
	}, nil
}

func (db *DB) OpenTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("BeginTx CreateCounterMetric failed, %w", err)
	}
	return tx, nil
}

func (db *DB) Rollback(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("rollback failed, %w", err)
	}
	return nil
}

func (db *DB) Commit(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit failed, %w", err)
	}
	return nil
}

func (db *DB) Shutdown(_ context.Context) {
	db.logger.Debug("db shutdown")
	db.pool.Close()
}
