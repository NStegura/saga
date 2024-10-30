package productcons

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Consumer        Consumer      `envconfig:"CONSUMER"`
	DB              DB            `envconfig:"DB"`
	Redis           Redis         `envconfig:"REDIS"`
	LogLevel        string        `envconfig:"LOG_LEVEL" default:"DEBUG"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
}

type Consumer struct {
	Topics  []string `envconfig:"TOPICS"  required:"true"`
	Brokers []string `envconfig:"BROKERS" required:"true"`
}

type DB struct {
	DSN string `envconfig:"DSN" required:"true"`
}

type Redis struct {
	DSN string `envconfig:"DSN" required:"true"`
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
