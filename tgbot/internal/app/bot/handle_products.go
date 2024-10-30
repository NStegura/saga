package bot

import (
	"errors"
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"

	"github.com/NStegura/saga/tgbot/internal/domain"
)

const (
	limit int64 = 5
)

func (b *TgBot) handleProducts() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		userState, err := b.cache.Get(b.ctx, ctx.Chat().ID)
		if err != nil {
			return fmt.Errorf("failed to get user state: %w", err)
		}

		switch ctx.Text() {
		case "/products":
			userState.ProductSt.CurPage = 0
		case "/products_next":
			userState.ProductSt.CurPage++
		case "/products_prev":
			if userState.ProductSt.CurPage > 0 {
				userState.ProductSt.CurPage--
			}
		default:
			return errors.New("unknown command")
		}
		b.logger.Info(userState.ProductSt.CurPage)
		products, err := b.business.GetProducts(b.ctx)
		if err != nil {
			return fmt.Errorf("failed to load products: %w", err)
		}
		b.logger.Infof("message handle: %s, %v", ctx.Text(), ctx.Chat())

		err = b.cache.Set(b.ctx, userState)
		if err != nil {
			return fmt.Errorf("failed to set user state: %w", err)
		}
		b.logger.Info(userState.ProductSt.CurPage)
		return ctx.Send(formatProducts(products, userState.ProductSt.CurPage, limit))
	}
}

func formatProducts(products []domain.Product, page, limit int64) string {
	// Todo: вынести пагинацию на сервер
	var isLastPage bool
	if len(products) == 0 {
		return "Нет доступных продуктов."
	}
	offset := page * limit
	end := offset + limit

	productsLen := int64(len(products))
	if productsLen < end {
		end = productsLen
	}

	if offset > end {
		offset = end
		isLastPage = true
	}

	products = products[offset:end]
	var sb strings.Builder
	sb.WriteString("Список продуктов:\n\n")

	for _, product := range products {
		sb.WriteString(fmt.Sprintf(
			"🛒 *%s* (Категория: %s)\n"+
				"📝 Описание: %s\n"+
				"📦 Количество: %d\n"+
				"🔖 Артикул: %d\n\n",
			product.Name, product.Category, product.Description, product.Count, product.ArticleID,
		))
	}
	sb.WriteString(fmt.Sprintf("Страница: %v\n", page+1))
	if isLastPage {
		sb.WriteString("⬅ /products_prev\n")
		sb.WriteString(fmt.Sprint("Больше нет товаров на продажу."))
	} else {
		sb.WriteString("⬅/products_prev⬅        ➡/products_next➡")
	}
	return sb.String()
}
