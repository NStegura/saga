package order

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/NStegura/saga/products/internal/app/consumers/handlers/order/models"
	"github.com/NStegura/saga/products/internal/clients/redis"
	productModels "github.com/NStegura/saga/products/internal/services/product/models"
	"github.com/sirupsen/logrus"
)

type IncomeHandler struct {
	product  Product
	cache    Cache
	orderCli OrderCli

	logger *logrus.Logger
}

func New(product Product, cache Cache, orderCli OrderCli, logger *logrus.Logger) *IncomeHandler {
	return &IncomeHandler{
		product:  product,
		cache:    cache,
		orderCli: orderCli,
		logger:   logger,
	}
}

func (i *IncomeHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }
func (i *IncomeHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	i.orderCli.Close()
	return nil
}

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
		productsToReserve, err := i.orderCli.GetProductsToReserve(message.OrderID)
		if err != nil {
			i.logger.Info("failed to reserve products: %w", err)
			continue
		}
		or := make([]productModels.Reserve, 0, len(productsToReserve))
		for _, product := range productsToReserve {
			or = append(or, productModels.Reserve{
				ProductID: product.ProductID,
				Count:     product.Count,
			})
		}
		err = i.product.ReserveProducts(ctx, message.OrderID, or)
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
