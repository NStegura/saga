package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	handler "github.com/NStegura/saga/orders/internal/app/consumers/handlers/payment"
	"github.com/NStegura/saga/orders/internal/services/order"

	"golang.org/x/sync/errgroup"

	config "github.com/NStegura/saga/orders/config/paymentcons"
	"github.com/NStegura/saga/orders/internal/app/consumers"
	"github.com/NStegura/saga/orders/internal/clients/kafka/consumer"
	"github.com/NStegura/saga/orders/internal/clients/redis"
	"github.com/NStegura/saga/orders/internal/storage"
	"github.com/NStegura/saga/orders/monitoring/logger"
)

const (
	groupID                      = "orders"
	inventoryConsumerServiceName = "payment consumer"
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

	orderService := order.New(repo, logg)

	cache := redis.New(cfg.Redis.DSN)

	cons := consumers.New(
		inventoryConsumerServiceName,
		cfg.Consumer.Topics,
		consGroup,
		handler.New(orderService, cache, logg),
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
