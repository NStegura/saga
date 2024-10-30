package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/NStegura/saga/products/internal/app/consumers"
	handler "github.com/NStegura/saga/products/internal/app/consumers/handlers/order"
	"github.com/NStegura/saga/products/internal/clients/orders"
	"github.com/NStegura/saga/products/internal/services/product"

	"golang.org/x/sync/errgroup"

	config "github.com/NStegura/saga/products/config/ordercons"
	"github.com/NStegura/saga/products/internal/clients/kafka/consumer"
	"github.com/NStegura/saga/products/internal/clients/redis"
	"github.com/NStegura/saga/products/internal/storage"
	"github.com/NStegura/saga/products/monitoring/logger"
)

const (
	groupID                      = "products"
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
	// для локального запуска необходимо подменить клиент
	orderCli, err := orders.New(cfg.OrderCli.CONN)
	if err != nil {
		return fmt.Errorf("failed to init cli: %w", err)
	}
	defer orderCli.Close()
	cons := consumers.New(
		inventoryConsumerServiceName,
		cfg.Consumer.Topics,
		consGroup,
		handler.New(productService, cache, orderCli, logg),
		logg,
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() (err error) {
		if err = cons.Start(ctx); err != nil {
			return fmt.Errorf("failed to start ctx: %w", err)
		}
		return
	})

	g.Go(func() (err error) {
		defer logg.Info("consumer has been shutdown")
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancelShutdownTimeoutCtx()
		err = cons.Shutdown(shutdownTimeoutCtx)
		if err != nil {
			return fmt.Errorf("failed to shutdown consumer: %w", err)
		}
		return
	})

	if err = g.Wait(); err != nil {
		return fmt.Errorf("failed to wait: %w", err)
	}
	return nil
}

func main() {
	if err := runConsumer(); err != nil {
		log.Fatal(err)
	}
}
