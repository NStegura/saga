package bot

import (
	"fmt"
	"github.com/NStegura/saga/tgbot/internal/domain"
	tele "gopkg.in/telebot.v4"
	"strconv"
	"strings"
)

func (b *TgBot) handleAddToShoppingCart() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}
		userState.ShopOrderSt.State = domain.WAIT_ADD_TO_SHOPPING_CART
		userState.ShopOrderSt.TryCount = 0
		err = b.cache.Set(b.ctx, userState)
		if err != nil {
			return fmt.Errorf("failed to set user state: %w", err)
		}
		return ctx.Send(
			"Чтобы добавить товары Вам необходимо прислать артикул и количество в формате: " +
				"Article,Count\n" +
				"Например: 1,10")
	}
}

func (b *TgBot) addToShoppingCart(msg string, userState *domain.UserState) (string, error) {
	parts := strings.Split(strings.TrimSpace(msg), ",")
	if len(parts) != 2 {
		return "", invalidFormatErr
	}
	productArticleID := strings.TrimSpace(parts[0])
	countStr := strings.TrimSpace(parts[1])

	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return "", invalidFormatErr
	}
	articleID, err := strconv.ParseInt(productArticleID, 10, 64)
	if err != nil {
		return "", invalidFormatErr
	}
	product, ok := userState.ShopCartSt.Products[articleID]
	if ok {
		userState.ShopCartSt.Products[articleID] = domain.OrderProduct{
			ArticleID: product.ArticleID,
			Count:     product.Count + count,
		}
	} else {
		userState.ShopCartSt.Products[articleID] = domain.OrderProduct{
			ArticleID: articleID,
			Count:     count,
		}
	}
	return "Товары были успешно добавлены в корзину.", nil
}
