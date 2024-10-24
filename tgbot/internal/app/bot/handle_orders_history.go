package bot

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"

	"github.com/NStegura/saga/tgbot/internal/domain"
)

func (b *TgBot) handleOrdersHistory() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}

		orders, err := b.business.GetOrders(b.ctx, userState.UserID)
		if err != nil {
			return fmt.Errorf("failed to get orders: %w", err)
		}

		return ctx.Send(formatOrders(orders))
	}
}

func formatOrders(orders []domain.Order) string {
	var sb strings.Builder

	if len(orders) == 0 {
		sb.WriteString("📭 У вас нет активных заказов.\n")
		return sb.String()
	}

	sb.WriteString("📋 *Список заказов:*\n\n")
	for _, order := range orders {
		// Выводим информацию о каждом заказе
		sb.WriteString(fmt.Sprintf("🧾 *Заказ #%d*\n", order.OrderID))
		sb.WriteString(fmt.Sprintf("🟢 Последний статус: %s\n", order.CurrentStatus))
		if len(order.StatusHistory) > 0 {
			sb.WriteString(
				fmt.Sprintf(
					"📅 Последнее обновление: %s\n",
					order.StatusHistory[len(order.StatusHistory)-1].Time.Format("2006-01-02 15:04:05")))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
