package easyredis

import "errors"

// ErrRedisKeyNotFound indicates that the cached element does not exist.
var ErrRedisKeyNotFound = errors.New("redis key not found")
