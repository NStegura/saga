package grpcserver

import "context"

type System interface {
	Ping(context.Context) error
}
