package consumers

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	name        string
	topics      []string
	consumerCli sarama.ConsumerGroup
	handler     sarama.ConsumerGroupHandler

	logger *logrus.Logger
}

func New(
	name string,
	topics []string,
	consumer sarama.ConsumerGroup,
	handler sarama.ConsumerGroupHandler,
	logger *logrus.Logger,
) *Consumer {
	return &Consumer{
		name:        name,
		topics:      topics,
		consumerCli: consumer,
		handler:     handler,
		logger:      logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		c.logger.Info("start consuming")
		if err := c.consumerCli.Consume(ctx, c.topics, c.handler); err != nil {
			c.logger.Errorf("income consumer error: %v", err)
			return fmt.Errorf("income consumer error: %w", err)
		}
		if ctx.Err() != nil {
			c.logger.Warningf("consumer ctx closed with err: %s", ctx.Err())
			return nil
		}
	}
}

func (c *Consumer) Name() string {
	return c.name
}

func (c *Consumer) Shutdown(ctx context.Context) (err error) {
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
