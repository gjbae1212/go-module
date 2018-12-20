package feedly

import (
	"github.com/joomcode/errorx"
)

var (
	feedlyError  = errorx.NewNamespace("[Err] feedly")
	emptyError   = feedlyError.NewType("empty")
	invalidError = feedlyError.NewType("invalid")
	unknownError = feedlyError.NewType("unknown")
)
