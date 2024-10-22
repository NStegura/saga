package product

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/orders/internal/app/consumers/handlers/product/models"
	"github.com/NStegura/saga/orders/internal/clients/redis"
	orderModels "github.com/NStegura/saga/orders/internal/services/order/models"
)

type IncomeHandler struct {
	order Order
	cache Cache

	logger *logrus.Logger
}

func New(order Order, cache Cache, logger *logrus.Logger) *IncomeHandler {
	return &IncomeHandler{
		order:  order,
		cache:  cache,
		logger: logger,
	}
}

func (i *IncomeHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (i *IncomeHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (i *IncomeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var (
			message models.ProductMessage
			state   orderModels.OrderState
		)

		i.logger.Debugf("Message:%s topic:%q partition:%d offset:%d\n",
			msg.Value, msg.Topic, msg.Partition, msg.Offset)

		err := json.Unmarshal(msg.Value, &message)
		if err != nil {
			i.logger.Infof("income data is valid %v: %v", string(msg.Value), err)
			continue
		}
		i.logger.Infof("Product event: %v", message)
		ctx := context.Background()

		switch message.Status {
		case models.FAILED:
			state = orderModels.RESERVE_FAILED
		case models.CREATED:
			state = orderModels.RESERVE_CREATED
		default:
			i.logger.Infof("product status: %s, continue", message.Status)
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

		err = i.order.CreateOrderState(ctx, message.OrderID, state)
		if err != nil {
			i.logger.Info("failed to update payment status: %w", err)
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