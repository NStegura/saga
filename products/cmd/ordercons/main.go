package main

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/products/internal/app/consumers"
	handler "github.com/NStegura/saga/products/internal/app/consumers/handlers/order"
	"github.com/NStegura/saga/products/internal/services/product"
	"log"
	"os"
	"os/signal"

	config "github.com/NStegura/saga/products/config/ordercons"
	"github.com/NStegura/saga/products/internal/clients/kafka/consumer"
	"github.com/NStegura/saga/products/internal/clients/redis"
	"github.com/NStegura/saga/products/internal/storage"
	"github.com/NStegura/saga/products/monitoring/logger"
	"golang.org/x/sync/errgroup"
)

const (
	groupID                      = "orders"
	inventoryConsumerServiceName = "order consumer"
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

	productService := product.New(repo, logg)

	cache := redis.New(cfg.Redis.DSN)

	cons := consumers.New(
		inventoryConsumerServiceName,
		cfg.Consumer.Topics,
		consGroup,
		handler.New(productService, cache, logg),
		logg,
	)

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

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancelShutdownTimeoutCtx()

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
