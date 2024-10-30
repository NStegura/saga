package bot

import (
	"fmt"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"

	"github.com/NStegura/saga/tgbot/internal/domain"
)

var answers = map[string]bool{
	"да":    true,
	"da":    true,
	"yes":   true,
	"true":  true,
	"нет":   false,
	"net":   false,
	"no":    false,
	"false": false,
}

func (b *TgBot) handlePay() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}
		userState.ShopOrderSt.State = domain.WAIT_PAY_ANSWER
		userState.ShopOrderSt.TryCount = 0
		err = b.cache.Set(b.ctx, userState)
		if err != nil {
			return fmt.Errorf("failed to set user state: %w", err)
		}
		return ctx.Send(
			"Чтобы оплатить зака, просто отправьте " +
				"'номер заказа, результат', например:\n" +
				"'1,да' для платы или '1,нет' для отмены.")
	}
}

func (b *TgBot) payOrder(msg string) (string, error) {
	parts := strings.Split(strings.TrimSpace(msg), ",")
	if len(parts) != 2 {
		return "", invalidFormatErr
	}
	oID := strings.TrimSpace(parts[0])
	agree := strings.TrimSpace(strings.ToLower(parts[1]))

	orderID, err := strconv.ParseInt(oID, 10, 64)
	if err != nil {
		return "", invalidFormatErr
	}

	res, ok := answers[agree]
	if !ok {
		return "", invalidFormatErr
	}

	err = b.business.PayOrder(b.ctx, orderID, res)
	if err != nil {
		return "", fmt.Errorf("failed to pay order: %w", err)
	}
	return "Платеж принят в обработку", nil
}
