package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	writeWait = 15 * time.Second

	pongWait = 60 * time.Second

	pingWait = 50 * time.Second
)

var (
	Upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
)

// It is a client between a server and a client
type Client struct {
	// breaker
	breaker Breaker

	// a connection of user
	conn *websocket.Conn

	// a buffer under send.
	send chan Message

	// message length
	maxMessageLength int64
}

func NewClient(bk *breaker, conn *websocket.Conn) (*Client, error) {
	client := &Client{
		breaker:          Breaker(bk),
		conn:             conn,
		send:             make(chan Message),
		maxMessageLength: bk.maxMessageLength,
	}
	go client.loopOfRead()
	go client.loopOfWrite()
	return client, nil
}

func (client *Client) loopOfRead() {
	defer func() {
		if err := client.breaker.UnRegister(client); err != nil {
			wraperr := errors.Wrap(err.(error), "[err] loopOfRead panic")
			client.breaker.(*breaker).errorHandler(wraperr)
		}
	}()

	client.conn.SetReadLimit(client.maxMessageLength)
	// it will be setting deadline of a websocket.
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	// it will be newly setting deadline of a websocket when a ping message received.
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		// If length of message a client is received will exceed over limit, it is raised error.
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// Message received will broadcast all of users.
		if err := client.breaker.BroadCast(&internalMessage{payload: message}); err != nil {
			wraperr := errors.Wrap(err.(error), "[err] loopOfRead broadcast")
			client.breaker.(*breaker).errorHandler(wraperr)
		}
	}
}

func (client *Client) loopOfWrite() {
	ticker := time.NewTicker(pingWait)
	defer func() {
		if err := client.breaker.UnRegister(client); err != nil {
			wraperr := errors.Wrap(err.(error), "[err] loopOfWrite panic")
			client.breaker.(*breaker).errorHandler(wraperr)
		}
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// If channel of the send is closed
			if !ok {
				if err := client.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					wraperr := errors.Wrap(err.(error), "[err] loopOfWrite close")
					client.breaker.(*breaker).errorHandler(wraperr)
				}
				return
			}
			// A message send to a user.
			if err := client.conn.WriteMessage(websocket.TextMessage, msg.GetMessage()); err != nil {
				wraperr := errors.Wrap(err.(error), "[err] loopOfWrite text")
				client.breaker.(*breaker).errorHandler(wraperr)
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// a ping message will periodically send.
			if err := client.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				wraperr := errors.Wrap(err.(error), "[err] loopOfWrite ping")
				client.breaker.(*breaker).errorHandler(wraperr)
				return
			}
		}
	}
}
