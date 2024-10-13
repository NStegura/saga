package order

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/NStegura/saga/products/internal/app/consumers/handlers/order/models"
	"github.com/NStegura/saga/products/internal/clients/redis"
	models2 "github.com/NStegura/saga/products/internal/services/product/models"
	"github.com/sirupsen/logrus"
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
		var message models.OrderMessage

		i.logger.Debugf("Message:%s topic:%q partition:%d offset:%d\n",
			msg.Value, msg.Topic, msg.Partition, msg.Offset)

		err := json.Unmarshal(msg.Value, &message)
		if err != nil {
			i.logger.Infof("income data is valid %v: %v", string(msg.Value), err)
			continue
		}
		i.logger.Infof("Order event: %v", message)
		ctx := context.Background()
		if message.Status != models.CREATED {
			i.logger.Infof("order status: %s, continue", message.Status)
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

		//добавить grpc клиента
		//ordersCli.GetOrder(orderID)

		productsInOrder := []models2.Reserve{
			{1, 5},
			{2, 1},
			{3, 3},
			{4, 5},
			{5, 2},
		}

		err = i.product.ReserveProducts(ctx, message.OrderID, productsInOrder)
		if err != nil {
			i.logger.Info("failed to reserve products: %w", err)
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
