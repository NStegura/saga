package main

import (
	"context"
	"github.com/NStegura/saga/golibs/event"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	ctx := context.Background()
	dsn := "postgres://usr:psswrd@localhost:54321/example?sslmode=disable"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	repo := event.NewEventRepository(pool, log)

}
