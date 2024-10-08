package main

import (
	"context"
	"fmt"
	config "github.com/NStegura/saga/payments/config/consumer"
	inventoryCons "github.com/NStegura/saga/payments/internal/app/consumer"
	"github.com/NStegura/saga/payments/internal/clients/kafka/consumer"
	"github.com/NStegura/saga/payments/internal/clients/redis"
	"github.com/NStegura/saga/payments/internal/services/payment"
	"github.com/NStegura/saga/payments/internal/storage"
	"github.com/NStegura/saga/payments/monitoring/logger"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
)

const (
	groupID                      = "payments"
	inventoryConsumerServiceName = "inventory consumer"
)

func runConsumer() error {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to init config: %w", err)
	}

	logg, err := logger.Init(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	consGroup, err := consumer.New(cfg.Consumer.Brokers, groupID, logg)
	if err != nil {
		return fmt.Errorf("failed to init consumerGroup")
	}

	repo, err := storage.New(ctx, cfg.DB.DSN, logg, false)
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	paymentService := payment.New(repo, logg)

	cache := redis.New(cfg.Redis.DSN)

	cons := inventoryCons.New(
		inventoryConsumerServiceName,
		cfg.Consumer.Topics,
		consGroup,
		paymentService,
		cache,
		logg)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() (err error) {
		if err = cons.Start(ctx); err != nil {
			return err
		}
		return
	})

	g.Go(func() (err error) {
		defer logg.Info("consumer has been shutdown")
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeutCtx := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancelShutdownTimeutCtx()

		return cons.Shutdown(shutdownTimeoutCtx)
	})

	if err = g.Wait(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := runConsumer(); err != nil {
		log.Fatal(err)
	}
}
