package bot

import (
	"fmt"
	tele "gopkg.in/telebot.v4"
	"strings"
)

func (b *TgBot) handleGetShoppingCart() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}

		if len(userState.ShopCartSt.Products) == 0 {
			return ctx.Send(
				"Корзина пуста, для ее заполнения:\n" +
					"1. Выберите продукты /products \n" +
					"2. Добавьте в корзину /add_to_shopping_cart\n",
			)
		}

		var sb strings.Builder
		sb.WriteString("Корзина товаров:\n\n")
		for articleID, shProduct := range userState.ShopCartSt.Products {
			product, err := b.business.GetProduct(b.ctx, articleID)
			if err != nil {
				return fmt.Errorf("failed to get product: %w", err)
			}
			sb.WriteString(fmt.Sprintf(
				"🛒 *%s* (Категория: %s)\n"+
					"📝 Описание: %s\n"+
					"📦 Количество: %d\n"+
					"🔖 Артикул: %d\n\n",
				product.Name, product.Category, product.Description, shProduct.Count, product.ArticleID,
			))
		}
		sb.WriteString(
			"Далее вы можете:\n" +
				"\t1. Cоздать заказ /create_order\n" +
				"\t2. Очистить корзину /clear_shopping_cart\n" +
				"\t3. Добавить товаров /add_to_shopping_cart")

		return ctx.Send(sb.String())
	}
}
