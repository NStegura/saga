package redis

import "errors"

var ErrCacheMiss = errors.New("not found data in cache")
