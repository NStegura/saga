package bot

import "errors"

var (
	invalidFormatErr = errors.New("неверный формат")
	orderNotFound    = errors.New("заказ не найден")
)
