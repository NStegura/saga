package server

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server Server `envconfig:"SERVER"`
	DB     DB     `envconfig:"DB"`
}

type Server struct {
	GRPCAddr string `envconfig:"GRPC_ADDR" default:"localhost:8081"`
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