package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"

	"github.com/NStegura/saga/tgbot/config"
	"github.com/NStegura/saga/tgbot/internal/app/bot"
	"github.com/NStegura/saga/tgbot/internal/clients/orders"
	"github.com/NStegura/saga/tgbot/internal/clients/payments"
	"github.com/NStegura/saga/tgbot/internal/clients/products"
	"github.com/NStegura/saga/tgbot/internal/clients/redis"
	"github.com/NStegura/saga/tgbot/internal/services/business"
	"github.com/NStegura/saga/tgbot/monitoring/logger"
)

func runTgBot() error {
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

	orderCli, err := orders.New(cfg.OrderCli.CONN, logg)
	if err != nil {
		return fmt.Errorf("failed to start order cli: %w", err)
	}

	paymentCli, err := payments.New(cfg.PaymentCli.CONN, logg)
	if err != nil {
		return fmt.Errorf("failed to start payment cli: %w", err)
	}

	productCli, err := products.New(cfg.ProductCli.CONN, logg)
	if err != nil {
		return fmt.Errorf("failed to start product cli: %w", err)
	}
	bus := business.New(orderCli, paymentCli, productCli, logg)

	cache := redis.New(cfg.Redis.DSN)
	tgBot, err := bot.New(cfg.TgBot.Token, cache, bus, logg)
	if err != nil {
		return fmt.Errorf("failed to start tg bot, %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() (err error) {
		logg.Info("start tgbot")
		if err = tgBot.Start(ctx); err != nil {
			return err
		}
		return
	})

	g.Go(func() (err error) {
		defer logg.Info("tgbot has been closed")
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeutCtx := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancelShutdownTimeutCtx()

		err = tgBot.Shutdown(shutdownTimeoutCtx)
		paymentCli.Close()
		orderCli.Close()
		productCli.Close()

		return
	})

	if err = g.Wait(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := runTgBot(); err != nil {
		log.Fatal(err)
	}
}
