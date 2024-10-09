package consumer

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func New(brokers []string, groupID string, logger *logrus.Logger) (sarama.ConsumerGroup, error) {
	logger.Info("start ConsumerGroup")
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to start ConsumerGroup: %w", err)
	}

	return group, nil
}
