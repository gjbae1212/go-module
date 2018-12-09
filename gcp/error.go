package gcp

import "github.com/joomcode/errorx"

var (
	GcpError     = errorx.NewNamespace("[Err] google cloud")
	EmptyError   = GcpError.NewType("empty")
	InvalidError = GcpError.NewType("invalid")
)
