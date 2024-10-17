package bot

import (
	"fmt"
	"github.com/NStegura/saga/tgbot/internal/domain"
	tele "gopkg.in/telebot.v4"
	"strconv"
	"strings"
)

func (b *TgBot) handleOrderInfo() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}
		userState.ShopOrderSt.State = domain.WAIT_ORDER_ID_TO_GET_ORDER_INFO
		userState.ShopOrderSt.TryCount = 0
		err = b.cache.Set(b.ctx, userState)
		if err != nil {
			return fmt.Errorf("failed to set user state: %w", err)
		}
		return ctx.Send("Чтобы получить заказ Вам необходимо передать его номер.")
	}
}

func (b *TgBot) getOrderInfo(msg string, userState *domain.UserState) (string, error) {
	orderID, err := strconv.ParseInt(msg, 10, 64)
	if err != nil {
		return "", invalidFormatErr
	}
	orders, err := b.business.GetOrders(b.ctx, userState.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to get orders: %w", err)
	}

	for _, order := range orders {
		if orderID == order.OrderID {
			return formatOrder(order), nil
		}
	}
	return "", orderNotFound
}

func formatOrder(order domain.Order) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("🧾 *Заказ #%d*\n", order.OrderID))
	sb.WriteString(fmt.Sprintf("📜 Описание: %s\n", order.Description))
	sb.WriteString(fmt.Sprintf("🟢 Статус: %s\n\n", order.CurrentStatus))

	sb.WriteString("📦 *Продукты:*\n")
	for _, product := range order.Products {
		sb.WriteString(fmt.Sprintf("  - *%s* (ID: %d)\n", product.Name, product.ArticleID))
		sb.WriteString(fmt.Sprintf("    📦 Категория: %s\n", product.Category))
		sb.WriteString(fmt.Sprintf("    📝 Описание: %s\n", product.Description))
		sb.WriteString(fmt.Sprintf("    🔢 Количество: %d\n", product.Count))
		sb.WriteString("\n")
	}

	sb.WriteString("📅 *История статусов:*\n")
	for _, status := range order.StatusHistory {
		sb.WriteString(
			fmt.Sprintf("  - %s (время: %s)\n", status.Status, status.Time.Format("2006-01-02 15:04:05")))
	}
	return sb.String()
}
