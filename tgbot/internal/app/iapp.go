package app

import "context"

type App interface {
	Start(context.Context) error
	Shutdown(context.Context) error
	Name() string
}
