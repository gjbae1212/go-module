package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type (
	// It is executed if error is raised
	ErrorHandler func(error)

	Breaker interface {
		Register(conn *websocket.Conn) error
		UnRegister(client *Client) error
		BroadCast(msg Message) error
	}

	// Breaker could manage all of something on the websocket.
	breaker struct {
		// ClientnMap for all of user connections.
		clientMap map[*Client]bool

		// Message for broadcast
		broadcast chan Message

		// Register connection
		register chan *Client

		// UnRegister connection
		unregister chan *Client

		// Error handler
		errorHandler ErrorHandler

		// Max read limit
		maxReadLimit int64
	}
)

// Create Breaker
func NewBreaker(opts ...Option) (Breaker, error) {
	bk := &breaker{
		clientMap:  make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	// default message length and error handler
	o := []Option{
		WithErrorHandlerOption(func(err error) {
			fmt.Printf("%+v \n", err)
		}),
		WithMaxReadLimit(512),
	}
	o = append(o, opts...)
	for _, opt := range o {
		opt.apply(bk)
	}
	// default message pool length
	if bk.broadcast == nil {
		WithMaxMessagePoolLength(100).apply(bk)
	}
	go bk.start()
	return Breaker(bk), nil
}

// WebSocket Connection is added to Breaker.
func (bk *breaker) Register(conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("[err] register empty params")
	}
	client, err := NewClient(bk, conn)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[err] register"))
	}
	bk.register <- client
	return nil
}

// Client is removed to Breaker
func (bk *breaker) UnRegister(client *Client) error {
	if client == nil {
		return fmt.Errorf("[err] unregister empty params")
	}
	bk.unregister <- client
	return nil
}

// Broadcast a message to all of clients.
func (bk *breaker) BroadCast(msg Message) error {
	if msg == nil {
		return fmt.Errorf("[err] message empty params")
	}
	bk.broadcast <- msg
	return nil
}

// Start a Breaker Loop
func (bk *breaker) start() {
	defer func() {
		if r := recover(); r != nil {
			wraperr := errors.Wrap(r.(error), "[err] breaker start panic")
			bk.errorHandler(wraperr)
			go bk.start()
		}
	}()

	for {
		select {
		case client := <-bk.register:
			bk.clientMap[client] = true
		case client := <-bk.unregister:
			if _, ok := bk.clientMap[client]; ok {
				delete(bk.clientMap, client)
				close(client.send)
				if err := client.conn.Close(); err != nil {
					wraperr := errors.Wrap(err, "[err] websocket close")
					bk.errorHandler(wraperr)
				}
			}
		case msg := <-bk.broadcast:
			for client, _ := range bk.clientMap {
				client.send <- msg
			}
		}
	}
}
