package system

import "github.com/sirupsen/logrus"

type System struct {
	repo   Repository
	logger *logrus.Logger
}

func New(repo Repository, logger *logrus.Logger) *System {
	return &System{repo: repo, logger: logger}
}
