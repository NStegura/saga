package server

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/golibs/event"
	config "github.com/NStegura/saga/payments/config/server"
	"github.com/NStegura/saga/payments/internal/app"
	"github.com/NStegura/saga/payments/internal/app/cron"
	"github.com/NStegura/saga/payments/internal/app/server"
	"github.com/NStegura/saga/payments/internal/clients/kafka/producer"
	"github.com/NStegura/saga/payments/internal/services/payment"
	"github.com/NStegura/saga/payments/internal/services/system"
	"github.com/NStegura/saga/payments/internal/storage"
	"github.com/NStegura/saga/payments/monitoring/logger"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"sync"
)

func runServer() error {
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

	repo, err := storage.New(ctx, cfg.DB.DSN, logg, true)
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	paymentService := payment.New(repo, logg)
	systemService := system.New(repo, logg)
	grpcServer, err := server.New(cfg.Server, paymentService, systemService, logg)
	if err != nil {
		return fmt.Errorf("failed to init grpc server: %w", err)
	}

	kafkaProducer, err := producer.New(cfg.Cron.Producer.Brokers, logg)
	if err != nil {
		return fmt.Errorf("failed to init kafka producer: %w", err)
	}
	worker := event.NewPushEventsWorker(
		kafkaProducer,
		event.NewEventService(repo, logg),
		logg,
		event.WithFrequency(cfg.Cron.Frequency),
		event.WithRateLimit(cfg.Cron.RateLimit),
		event.WithEventsLimit(cfg.Cron.EventsLimit),
		event.WithReserve(cfg.Cron.Reserve),
	)

	cronTab := cron.New(worker, logg)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		if err = grpcServer.Start(ctx); err != nil {
			return err
		}
		return
	})

	g.Go(func() (err error) {
		if err = cronTab.Start(ctx); err != nil {
			return err
		}
		return
	})

	g.Go(func() (err error) {
		defer logg.Info("server has been shutdown")
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeutCtx := context.WithTimeout(
			context.Background(), cfg.Server.ShutdownTimeout)
		defer cancelShutdownTimeutCtx()

		var wg sync.WaitGroup
		processes := []app.App{cronTab, grpcServer}
		wg.Add(len(processes))

		var errors []error
		for _, proc := range processes {
			go func(proc app.App) {
				defer wg.Done()
				if err = proc.Shutdown(shutdownTimeoutCtx); err != nil {
					logg.Errorf("%s shutdown error: %v", proc.Name(), err)
				}
			}(proc)
		}
		wg.Wait()
		if len(errors) > 0 {
			return fmt.Errorf("shutdown encountered errors: %v", errors)
		}
		return nil
	})

	if err = g.Wait(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := runServer(); err != nil {
		log.Fatal(err)
	}
}
