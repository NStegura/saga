package bot

import (
	"errors"
	"fmt"
	"github.com/NStegura/saga/tgbot/internal/domain"
	tele "gopkg.in/telebot.v4"
	"strconv"
	"strings"
)

func (b *TgBot) handleAddToOrder() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}
		userState.ShopOrderSt.State = domain.WAIT_ADD_ORDER
		userState.ShopOrderSt.TryCount = 0
		err = b.cache.Set(b.ctx, userState)
		if err != nil {
			return fmt.Errorf("failed to set user state: %w", err)
		}
		return ctx.Send(
			"Чтобы добавить товары вам необходимо прислать артикул и количество в формате: " +
				"Article,Count\n" +
				"Например: 1,10")
	}
}

func addToOrder(msg string, userState *domain.UserState) error {
	parts := strings.Split(strings.TrimSpace(msg), ",")
	if len(parts) != 2 {
		return errors.New("invalid format")
	}
	productArticleID := strings.TrimSpace(parts[0])
	countStr := strings.TrimSpace(parts[1])

	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return errors.New("invalid format")
	}
	articleID, err := strconv.ParseInt(productArticleID, 10, 64)
	if err != nil {
		return errors.New("invalid format")
	}
	product, ok := userState.ProductSt.Products[articleID]
	if ok {
		userState.ProductSt.Products[articleID] = domain.OrderProduct{
			ArticleID: product.ArticleID,
			Count:     product.Count + count,
		}
	} else {
		orderProductMap := make(map[int64]domain.OrderProduct, 10)
		orderProductMap[articleID] = domain.OrderProduct{
			ArticleID: articleID,
			Count:     count,
		}
		userState.ProductSt.Products = orderProductMap
	}
	return nil
}
