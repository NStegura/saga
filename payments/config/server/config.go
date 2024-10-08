package server

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Server   Server `envconfig:"SERVER"`
	Cron     Cron   `envconfig:"CRON"`
	DB       DB     `envconfig:"DB"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"DEBUG"`
}

type Cron struct {
	Producer    Producer      `envconfig:"PRODUCER"`
	Frequency   time.Duration `envconfig:"FREQUENCY" default:"5s"`
	RateLimit   int           `envconfig:"RATE_LIMIT" default:"2"`
	EventsLimit int           `envconfig:"EVENTS_LIMIT" default:"20"`
	Reserve     time.Duration `envconfig:"RESERVE" default:"20s"`
}

type Producer struct {
	Brokers []string `envconfig:"BROKERS" required:"true"`
}

type Server struct {
	GRPCAddr        string        `envconfig:"GRPC_ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
}

type DB struct {
	DSN string `envconfig:"DSN"  required:"true"`
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
