package payment

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/products/internal/app/consumers/handlers/payment/models"
	"github.com/NStegura/saga/products/internal/clients/redis"
)

type IncomeHandler struct {
	product Product
	cache   Cache

	logger *logrus.Logger
}

func New(product Product, cache Cache, logger *logrus.Logger) *IncomeHandler {
	return &IncomeHandler{
		product: product,
		cache:   cache,
		logger:  logger,
	}
}

func (i *IncomeHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (i *IncomeHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (i *IncomeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var (
			message models.PaymentMessage
			status  bool
		)

		i.logger.Debugf("Message:%s topic:%q partition:%d offset:%d\n",
			msg.Value, msg.Topic, msg.Partition, msg.Offset)

		err := json.Unmarshal(msg.Value, &message)
		if err != nil {
			i.logger.Infof("income data is valid %v: %v", string(msg.Value), err)
			continue
		}
		i.logger.Infof("Payment event: %v", message)
		ctx := context.Background()

		switch message.Status {
		case models.FAILED:
			status = false
		case models.COMPLETED:
			status = true
		default:
			i.logger.Infof("unknown payment status: %s, continue", message.Status)
			continue
		}

		err = i.cache.Get(ctx, message.IKey)
		if err == nil {
			i.logger.Info("idempotent key already exists, continue")
			continue
		}

		if !errors.Is(err, redis.ErrCacheMiss) {
			i.logger.Infof("unexpected err: %v", err)
			continue
		}

		err = i.product.UpdateReserveStatus(ctx, message.OrderID, status)
		if err != nil {
			i.logger.Info("failed to update reserve status: %w", err)
			continue
		}
		err = i.cache.Set(ctx, message.IKey)
		if err != nil {
			i.logger.Info("failed to set key: %w", err)
			continue
		}
		session.MarkMessage(msg, "")
	}

	return nil
}
