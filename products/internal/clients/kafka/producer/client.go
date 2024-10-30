package producer

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	producer sarama.SyncProducer
	logger   *logrus.Logger
}

func New(brokers []string, logger *logrus.Logger) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to sync producer: %w", err)
	}
	return &Producer{
		producer: syncProducer,
		logger:   logger,
	}, nil
}

func (p *Producer) PushMsg(msg []byte, topic string) error {
	par, off, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("sync"),
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		return fmt.Errorf("failed to push msg %v", msg)
	}
	p.logger.Debugf("order %v -> %v; %v", par, off, err)
	return nil
}

func (p *Producer) PushMsgs(msgs [][]byte, topic string) error {
	pms := make([]*sarama.ProducerMessage, 0, len(msgs))
	for _, msg := range msgs {
		pms = append(pms, &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder("sync"),
			Value: sarama.ByteEncoder(msg),
		})
	}
	if err := p.producer.SendMessages(pms); err != nil {
		return fmt.Errorf("failed to push msgs %v", msgs)
	}

	p.logger.Debug("pushed msgs to topic payment")
	return nil
}
