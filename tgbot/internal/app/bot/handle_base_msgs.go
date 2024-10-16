package bot

import (
	"fmt"
	"github.com/NStegura/saga/tgbot/internal/domain"
	tele "gopkg.in/telebot.v4"
)

const maxTryCountToReadMsg = 3

func (b *TgBot) baseHandle() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		// Todo: вынести userState в middleware с получением и сохранением, передача везде по указателю из ctx
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}

		if userState.ShopOrderSt.TryCount > maxTryCountToReadMsg {
			return ctx.Send("Количество попыток для команды исчерпано, попробуйте еще раз")
		}
		switch userState.ShopOrderSt.State {
		case domain.WAIT_ADD_ORDER:
			err = addToOrder(ctx.Text(), &userState)
			if err != nil {
				userState.ShopOrderSt.TryCount++
			} else {
				userState.ShopOrderSt.State = domain.NONE
			}
		case domain.NONE:
			return ctx.Send("Кажется, что вы пишите просто так, попробуйте /start")
		default:
			return ctx.Send("Кажется, что вы пишите просто так, попробуйте /start")
		}
		if err != nil {
			return ctx.Send(err.Error())
		}

		err = b.cache.Set(b.ctx, userState)
		if err != nil {
			return fmt.Errorf("failed to set user state: %w", err)
		}
		return ctx.Send("Товары были успешно добавлены в корзину.")
	}
}
