package payment

import (
	"context"
	"encoding/json"
	"fmt"
	dbModels "github.com/NStegura/saga/payments/internal/repo/models"
	"github.com/NStegura/saga/payments/internal/services/payment/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	topic = "payment"
)

type Payment struct {
	repo      Repository
	eventRepo EventRepository
	logger    *logrus.Logger
}

func New(repo Repository, logger *logrus.Logger) *Payment {
	return &Payment{repo: repo, logger: logger}
}

func (p *Payment) CreatePayment(ctx context.Context, orderID int64) (id int64, err error) {
	tx, err := p.repo.OpenTransaction(ctx)
	if err != nil {
		return id, fmt.Errorf("failed to open transaction, %w", err)
	}
	defer func() {
		_ = p.repo.Commit(ctx, tx)
	}()

	payID, err := p.repo.CreatePayment(ctx, tx, orderID)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return 0, fmt.Errorf("failed to create payment, %w", err)
	}

	paymentMessage := models.PaymentMessage{
		IKey:    uuid.New(),
		OrderID: orderID,
		Status:  models.CREATED,
	}
	payload, err := json.Marshal(paymentMessage)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return id, fmt.Errorf("failed to marshal payment message: %w", err)
	}

	err = p.eventRepo.CreateEvent(ctx, tx, topic, payload)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return id, fmt.Errorf("failed to create event: %w", err)
	}
	return payID, err
}

func (p *Payment) UpdatePaymentStatus(ctx context.Context, orderID int64, status models.PaymentMessageStatus) (err error) {
	var dbStatus dbModels.PaymentStatus

	tx, err := p.repo.OpenTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to open transaction, %w", err)
	}
	defer func() {
		_ = p.repo.Commit(ctx, tx)
	}()

	_, err = p.repo.GetPaymentByOrderID(ctx, tx, orderID, true)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to get payment by order id for update: %w", err)
	}

	if status == models.COMPLETED {
		dbStatus = dbModels.COMPLETED
	} else if status == models.FAILED {
		dbStatus = dbModels.FAILED
	} else {
		_ = p.repo.Rollback(ctx, tx)
		return fmt.Errorf("unknown status")
	}

	err = p.repo.UpdatePaymentStatusByOrderID(ctx, tx, orderID, dbStatus)
	if err != nil {
		p.logger.Errorf("update status err: %v", err)
		_ = p.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to get payment by order id for update: %w", err)
	}

	paymentMessage := models.PaymentMessage{
		IKey:    uuid.New(),
		OrderID: orderID,
		Status:  status,
	}
	payload, err := json.Marshal(paymentMessage)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to marshal payment message: %w", err)
	}
	err = p.eventRepo.CreateEvent(ctx, tx, topic, payload)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to create event: %w", err)
	}
	return
}
