package producer

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	topic    string
	producer sarama.SyncProducer
	logger   *logrus.Logger
}

func New(brokers []string, topic string, logger *logrus.Logger) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		logger.Fatalf("sync kafka: %v", err)
	}
	return &Producer{
		topic:    topic,
		producer: syncProducer,
		logger:   logger,
	}, nil
}

func (p *Producer) PushMsg(msg []byte) error {
	par, off, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic, //payment
		Key:   sarama.StringEncoder("sync"),
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		return fmt.Errorf("failed to push msg %v", msg)
	}
	p.logger.Debug("order %v -> %v; %v", par, off, err)
	return nil
}

func (p *Producer) PushMsgs(msgs [][]byte) error {
	pms := make([]*sarama.ProducerMessage, len(msgs))
	for _, msg := range msgs {
		pms = append(pms, &sarama.ProducerMessage{
			Topic: p.topic, //payment
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
