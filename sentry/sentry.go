package sentry

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"net/http"
	"runtime"
)

func Load(sentryDSN string) error {
	if err := raven.SetDSN(sentryDSN); err != nil {
		return err
	}
	return nil
}

func RaiseSentry(err error, stack bool) {
	message := fmt.Sprint(err)
	packet := raven.NewPacket(message)
	if stack {
		trace := make([]byte, 1024*50)
		runtime.Stack(trace, false)
		packet.Extra["Stack"] = string(trace)
	}
	raven.Capture(packet, nil)
}

func RaiseSentryWithRequest(err error, req *http.Request, stack bool) {
	message := fmt.Sprint(err)
	packet := raven.NewPacket(message, raven.NewHttp(req))
	if stack {
		trace := make([]byte, 1024*50)
		runtime.Stack(trace, false)
		packet.Extra["Stack"] = string(trace)
	}
	raven.Capture(packet, nil)
}
