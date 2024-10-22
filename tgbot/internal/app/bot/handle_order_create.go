package bot

import (
	"fmt"

	tele "gopkg.in/telebot.v4"

	"github.com/NStegura/saga/tgbot/internal/domain"
)

func (b *TgBot) handleCreateOrder() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}

		if len(userState.ShopCartSt.Products) == 0 {
			return ctx.Send("Товаров в корзине нет.")
		}

		ops := make([]domain.OrderProduct, 0, len(userState.ShopCartSt.Products))
		for _, orderProduct := range userState.ShopCartSt.Products {
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

		userState.ShopCartSt.Products = map[int64]domain.OrderProduct{}
		err = b.cache.Set(b.ctx, userState)
		if err != nil {
			return fmt.Errorf("failed to set user state: %w", err)
		}
		return ctx.Send(fmt.Sprintf("Заказ с номером %d был создан", orderID))
	}
}
