package server

type Config struct {
	Server Server `envconfig:"SERVER"`
	DB     DB
}

type Server struct {
}

type DB struct {
}
