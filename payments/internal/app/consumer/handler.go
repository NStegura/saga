package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/payments/internal/app/consumer/models"
	"github.com/NStegura/saga/payments/internal/clients/redis"
)

type IncomeHandler struct {
	payment Payments
	cache   Cache

	logger *logrus.Logger
}

func (i *IncomeHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (i *IncomeHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (i *IncomeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var message models.InventoryMessage
		i.logger.Debugf("Message:%v topic:%q partition:%d offset:%d\n",
			msg.Value, msg.Topic, msg.Partition, msg.Offset)

		err := json.Unmarshal(msg.Value, &message)
		if err != nil {
			i.logger.Infof("income data is valid %v: %v", string(msg.Value), err)
			continue
		}
		i.logger.Infof("Inventory event: %v", message)
		ctx := context.Background()
		if message.Status != models.COMPLETED {
			i.logger.Infof("inventory status: %s, continue", message.Status)
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

		_, err = i.payment.CreatePayment(ctx, message.OrderID)
		if err != nil {
			i.logger.Info("failed to create payment: %w", err)
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
