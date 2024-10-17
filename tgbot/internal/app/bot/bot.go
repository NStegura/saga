package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v4"
)

const (
	botName = "tg_order_bot"
)

type TgBot struct {
	ctx      context.Context
	bot      *tele.Bot
	cache    Cache
	business Business
	name     string

	logger *logrus.Logger
}

func New(token string, cache Cache, business Business, logger *logrus.Logger) (*TgBot, error) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 1 * time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init tgbot: %w", err)
	}
	return &TgBot{
		bot:      bot,
		name:     botName,
		cache:    cache,
		business: business,
		logger:   logger,
	}, nil
}

func (b *TgBot) configRoutes(ctx context.Context) {
	b.ctx = ctx
	b.bot.Handle(tele.OnText, b.baseHandle())
	b.bot.Handle("/start", b.handleStart())
	b.bot.Handle("/products", b.handleProducts())
	b.bot.Handle("/products_next", b.handleProducts())
	b.bot.Handle("/products_prev", b.handleProducts())

	b.bot.Handle("/get_shopping_cart", b.handleGetShoppingCart())
	b.bot.Handle("/add_to_shopping_cart", b.handleAddToShoppingCart())
	b.bot.Handle("/clear_shopping_cart", b.handleClearShoppingCart())

	b.bot.Handle("/order_info", b.handleOrderInfo())
	b.bot.Handle("/orders_history", b.handleOrdersHistory())
	b.bot.Handle("/create_order", b.handleCreateOrder())

	b.bot.Handle("/pay", b.handlePay())

	b.bot.Handle("/help", b.handleHelp())
}

func (b *TgBot) Start(ctx context.Context) error {
	b.configRoutes(ctx)
	b.bot.Start()
	return nil
}

func (b *TgBot) Shutdown(_ context.Context) error {
	b.bot.Stop()
	return nil
}

func (b *TgBot) Name() string {
	return b.name
}
