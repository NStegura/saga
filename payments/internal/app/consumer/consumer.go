package consumer

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type InventoryConsumer struct {
	name        string
	ctx         context.Context
	topics      []string
	consumerCli sarama.ConsumerGroup
	payment     Payments
	cache       Cache

	logger *logrus.Logger
}

func New(
	name string,
	topics []string,
	consumer sarama.ConsumerGroup,
	payment Payments,
	cache Cache,
	logger *logrus.Logger,
) *InventoryConsumer {
	return &InventoryConsumer{
		name:        name,
		topics:      topics,
		consumerCli: consumer,
		payment:     payment,
		cache:       cache,
		logger:      logger,
	}
}

func (c *InventoryConsumer) Start(ctx context.Context) error {
	for {
		c.logger.Info("start consuming")
		if err := c.consumerCli.Consume(ctx, c.topics, &IncomeHandler{
			payment: c.payment,
			cache:   c.cache,
		}); err != nil {
			c.logger.Errorf("income consumer error: %v", err)
			return fmt.Errorf("income consumer error: %w", err)
		}
		if ctx.Err() != nil {
			c.logger.Warning("consumer ctx closed with err: %v", ctx.Err())
			return nil
		}
	}
}

func (c *InventoryConsumer) Name() string {
	return c.name
}

func (c *InventoryConsumer) Shutdown(ctx context.Context) (err error) {
	doneCh := make(chan struct{})
	errCh := make(chan error)

	go func() {
		defer close(doneCh)
		defer close(errCh)

		if err = c.consumerCli.Close(); err != nil {
			errCh <- fmt.Errorf("shutdown consumer failed")
			return
		}
		doneCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
		c.logger.Info("Shutdown timeout reached, force closing.")
	case err = <-errCh:
		c.logger.Errorf("Shutdown consumer with err: %v", err)
	case <-doneCh:
		c.logger.Info("Shutdown success")
	}
	return
}
