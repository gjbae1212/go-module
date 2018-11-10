package redis

import (
	"github.com/joomcode/errorx"
)

var (
	redisError = errorx.NewNamespace("[Err] redis")

	emptyError    = redisError.NewType("empty")
	notFoundError = redisError.NewType("not_found")
	unknownError  = redisError.NewType("unknown")
	invalidError  = redisError.NewType("invalid")
	alreadyError  = redisError.NewType("already")
)
