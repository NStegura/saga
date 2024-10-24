package server

import "context"

type System interface {
	Ping(context.Context) error
}
