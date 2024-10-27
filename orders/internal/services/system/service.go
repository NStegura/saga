package system

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type System struct {
	repo   Repository
	logger *logrus.Logger
}

func New(repo Repository, logger *logrus.Logger) *System {
	return &System{repo: repo, logger: logger}
}

func (s *System) Ping(ctx context.Context) error {
	if err := s.repo.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}
	return nil
}
