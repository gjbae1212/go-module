package sentry

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"net/http"
	"runtime"
)

const (
	defaultStackLength = 20 * 1024
)

var (
	stackLength int
)

func Load(sentryDSN string, slength int) error {
	if err := raven.SetDSN(sentryDSN); err != nil {
		return err
	}
	if slength <= 0 {
		stackLength = defaultStackLength
	} else {
		stackLength = slength
	}

	return nil
}

func MakePacket(err error, stack bool) *raven.Packet {
	message := fmt.Sprint(err)
	packet := raven.NewPacket(message)
	if stack {
		trace := make([]byte, stackLength)
		runtime.Stack(trace, false)
		packet.Extra["Stack"] = string(trace)
	}
	return packet
}

func MakePacketWithRequest(err error, req *http.Request, stack bool) *raven.Packet {
	message := fmt.Sprint(err)
	packet := raven.NewPacket(message, raven.NewHttp(req))
	if stack {
		trace := make([]byte, stackLength)
		runtime.Stack(trace, false)
		packet.Extra["Stack"] = string(trace)
	}
	return packet
}

func Raise(packet *raven.Packet, captureTags map[string]string) {
	raven.Capture(packet, captureTags)
}
