package bot

import (
	"errors"
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"

	"github.com/NStegura/saga/tgbot/internal/domain"
)

const maxTryCountToReadMsg = 3

func (b *TgBot) baseHandle() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		var msg string
		// Todo: вынести userState в middleware с получением и сохранением, передача везде по указателю из ctx
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}

		if userState.ShopOrderSt.TryCount > maxTryCountToReadMsg {
			return ctx.Send("Количество попыток для команды исчерпано, попробуйте еще раз")
		}
		switch userState.ShopOrderSt.State {
		case domain.WAIT_ADD_TO_SHOPPING_CART:
			msg, err = b.addToShoppingCart(ctx.Text(), &userState)
		case domain.WAIT_ORDER_ID_TO_GET_ORDER_INFO:
			msg, err = b.getOrderInfo(ctx.Text(), &userState)
		case domain.WAIT_PAY_ANSWER:
			msg, err = b.payOrder(ctx.Text())
		case domain.NONE:
			msg = "Кажется, что вы пишите просто так, попробуйте /start"
		default:
			msg = "Кажется, что вы пишите просто так, попробуйте /start"
		}
		if err != nil {
			var sb strings.Builder
			userState.ShopOrderSt.TryCount++
			if errors.Is(invalidFormatErr, err) {
				sb.WriteString(fmt.Sprintf("Ошибка: %s\n", err.Error()))
			} else {
				b.logger.Info(err.Error())
				sb.WriteString(fmt.Sprintf("Что-то пошло не так\n"))
			}
			sb.WriteString(fmt.Sprintf(
				"Кол-во попыток осталось: %d",
				maxTryCountToReadMsg-userState.ShopOrderSt.TryCount,
			))
			return ctx.Send(sb)
		}
		userState.ShopOrderSt.State = domain.NONE

		err = b.cache.Set(b.ctx, userState)
		if err != nil {
			return fmt.Errorf("failed to set user state: %w", err)
		}
		return ctx.Send(msg)
	}
}
