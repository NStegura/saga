package config

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	TgBot           TgBot         `envconfig:"TG_BOT"`
	Redis           Redis         `envconfig:"REDIS"`
	OrderCli        OrderCli      `envconfig:"ORDER_CLI"`
	PaymentCli      PaymentCli    `envconfig:"PAYMENT_CLI"`
	ProductCli      ProductCli    `envconfig:"PRODUCT_CLI"`
	LogLevel        string        `envconfig:"LOG_LEVEL" default:"DEBUG"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
}

type TgBot struct {
	Token string `envconfig:"TOKEN" required:"true" default:"7588716328:AAFpRkPVV_34C0_tahmCrOuQSKSQYudqEG8"` // Todo: put away default(Debug)
}

type Redis struct {
	DSN string `envconfig:"DSN" required:"true" default:"0.0.0.0:6379"` // Todo: put away default(Debug)
}

type OrderCli struct {
	CONN string `envconfig:"CONN" required:"true" default:"0.0.0.0:8081"` // Todo: put away default(Debug)
}

type PaymentCli struct {
	CONN string `envconfig:"CONN" required:"true" default:"0.0.0.0:8083"` // Todo: put away default(Debug)
}

type ProductCli struct {
	CONN string `envconfig:"CONN" required:"true" default:"0.0.0.0:8082"` // Todo: put away default(Debug)
}

func New() (cfg Config, err error) {
	if err = envconfig.Process("", &cfg); err != nil {
		return
	}
	if err = yaml.NewEncoder(os.Stdout).Encode(&cfg); err != nil {
		return
	}
	return
}
