package cron

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"

	mock_workers "github.com/NStegura/saga/payments/mocks/app/cron/workers"
)

func TestCron_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWorker := mock_workers.NewMockWorker(ctrl)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	mockWorker.EXPECT().GetFrequency().Return(time.Millisecond * 10).AnyTimes()
	mockWorker.EXPECT().Run(gomock.Any()).Return(nil).Times(6)

	cronJob := New(mockWorker, logger)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*60)
	defer cancel()

	go func() {
		err := cronJob.Start(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	time.Sleep(time.Millisecond * 70)
}

func TestCron_WorkerFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWorker := mock_workers.NewMockWorker(ctrl)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	mockWorker.EXPECT().GetFrequency().Return(time.Millisecond * 10).AnyTimes()
	mockWorker.EXPECT().Run(gomock.Any()).Return(assert.AnError).Times(1)
	mockWorker.EXPECT().Run(gomock.Any()).Return(nil).Times(5)

	cronJob := New(mockWorker, logger)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*60)
	defer cancel()

	go func() {
		err := cronJob.Start(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	time.Sleep(time.Millisecond * 70)
}
