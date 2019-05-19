package websocket

// It is interface on message.
type (
	Message interface {
		GetMessage() []byte
	}

	internalMessage struct {
		payload []byte
	}
)

func (im *internalMessage) GetMessage() []byte {
	return im.payload
}
