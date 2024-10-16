package bot

import (
	"fmt"
	"github.com/NStegura/saga/tgbot/internal/domain"
	tele "gopkg.in/telebot.v4"
)

func (b *TgBot) handleCreateOrder() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}

		if len(userState.ProductSt.Products) == 0 {
			return ctx.Send("Товаров в корзине нет.")
		}

		ops := make([]domain.OrderProduct, 0, len(userState.ProductSt.Products))
		for _, orderProduct := range userState.ProductSt.Products {
			ops = append(ops, orderProduct)
		}

		orderID, err := b.business.CreateOrder(b.ctx, domain.CreateOrderInfo{
			UserID:      userState.UserID,
			Description: "",
			Products:    ops,
		})
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		userState.ProductSt.Products = map[int64]domain.OrderProduct{}
		err = b.cache.Set(b.ctx, userState)
		if err != nil {
			return fmt.Errorf("failed to set user state: %w", err)
		}
		return ctx.Send(fmt.Sprintf("Заказ с номером %d был создан", orderID))
	}
}
