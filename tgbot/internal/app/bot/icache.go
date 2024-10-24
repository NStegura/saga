package bot

import (
	"context"

	"github.com/NStegura/saga/tgbot/internal/domain"
)

type Cache interface {
	Get(ctx context.Context, userID int64) (userState domain.UserState, err error)
	Set(ctx context.Context, userState domain.UserState) (err error)
	Del(ctx context.Context, userID int64) error
}
