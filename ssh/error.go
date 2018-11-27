package ssh

import "github.com/joomcode/errorx"

var (
	SSHError           = errorx.NewNamespace("[Err] ssh")
	EmptyError         = SSHError.NewType("empty")
	InvalidParamsError = SSHError.NewType("invalid params")
)
