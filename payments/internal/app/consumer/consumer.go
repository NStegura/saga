package consumer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/NStegura/saga/payments/internal/services/payment"
	"github.com/sirupsen/logrus"
	"time"
)

type InventoryConsumer struct {
	ctx         context.Context
	topics      []string
	consumerCli sarama.ConsumerGroup
	payment     payment.Payment

	logger *logrus.Logger
}

func New(
	topics []string,
	consumer sarama.ConsumerGroup,
	payment payment.Payment,
	logger *logrus.Logger,
) *InventoryConsumer {
	return &InventoryConsumer{
		topics:      topics,
		consumerCli: consumer,
		payment:     payment,
		logger:      logger,
	}
}

func (c *InventoryConsumer) Start(ctx context.Context) {
	for {
		c.logger.Info("start consuming")
		handler := &IncomeHandler{payment: c.payment}
		err := c.consumerCli.Consume(ctx, c.topics, handler)
		if err != nil {
			c.logger.Errorf("income consumer error: %v", err)
			time.Sleep(time.Second * 5)
		}
	}
}
