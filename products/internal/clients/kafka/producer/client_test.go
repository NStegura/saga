package producer

import (
	"testing"

	"github.com/IBM/sarama/mocks"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPushMsgSuccess(t *testing.T) {
	logger := logrus.New()

	mockProducer := mocks.NewSyncProducer(t, nil)
	mockProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
		if string(val) != "test message" {
			t.Errorf("expected test message, got %v", val)
		}
		return nil
	})

	p := &Producer{
		producer: mockProducer,
		logger:   logger,
	}

	msg := []byte("test message")
	topic := "test-topic"

	err := p.PushMsg(msg, topic)
	assert.NoError(t, err)
}

func TestPushMsgsSuccess(t *testing.T) {
	logger := logrus.New()

	mockProducer := mocks.NewSyncProducer(t, nil)

	mockProducer.ExpectSendMessageAndSucceed()
	mockProducer.ExpectSendMessageAndSucceed()

	p := &Producer{
		producer: mockProducer,
		logger:   logger,
	}

	msgs := [][]byte{
		[]byte("test message 1"),
		[]byte("test message 2"),
	}
	topic := "test-topic"

	err := p.PushMsgs(msgs, topic)
	assert.NoError(t, err)
}
