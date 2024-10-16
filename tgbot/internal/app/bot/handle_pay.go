package bot

import tele "gopkg.in/telebot.v4"

func (b *TgBot) handlePay() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		b.logger.Infof("message handle: %s", ctx.Text())
		return ctx.Send("Добро пожаловать! Используйте /create_order для создания заказа.")
	}
}
