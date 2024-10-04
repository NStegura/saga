package cron

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	PushJob PushJob `envconfig:"CRON"`
	DB      DB      `envconfig:"DB"`
}

type PushJob struct {
	Frequency   time.Duration `envconfig:"FREQUENCY" default:"10s"`
	RateLimit   int           `envconfig:"RATE_LIMIT" default:"2"`
	EventsLimit int           `envconfig:"EVENTS_LIMIT" default:"20"`
}

type DB struct {
	DSN string `envconfig:"DSN" default:"postgres://usr:psswrd@localhost:54321/payments?sslmode=disable"`
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
