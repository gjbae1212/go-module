package websocket

import (
	"testing"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var (
	mockBreaker Breaker
)

func init() {
	mockBreaker, _ = NewBreaker(WithMaxReadLimit(10), WithMaxMessagePoolLength(1))
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"

	var conns []*websocket.Conn
	for i := 0; i < 10; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		assert.NoError(err)
		conns = append(conns, conn)
		go func(conn *websocket.Conn) {
			for {
				_, _, err := conn.ReadMessage() // if a ping message is received, client will send pong.
				if err != nil {
					return
				}
			}
		}(conn)
	}
	time.Sleep(1 * time.Second)
	assert.Len(mockBreaker.(*breaker).clientMap, 10)

	// BROADCAST
	for i := 0; i < 10; i++ {
		err := mockBreaker.BroadCast(&internalMessage{payload: []byte("HELLO WORLD ABCDEFGHIJKNMKOPKQSTUVWXZY")})
		assert.NoError(err)
	}

	// PING, PONG CHECK
	//time.Sleep(70 * time.Second)
	//assert.Len(mockBreaker.(*breaker).clientMap, 10)

	// READLIMIT
	for _, conn := range conns {
		err := conn.WriteMessage(websocket.TextMessage, []byte("check"))
		assert.NoError(err)
	}
	time.Sleep(1 * time.Second)
	assert.Len(mockBreaker.(*breaker).clientMap, 10)
	for _, conn := range conns {
		err := conn.WriteMessage(websocket.TextMessage, []byte("raise error!!!!!!!!!!!!!!!!!"))
		assert.NoError(err)
	}
	time.Sleep(1 * time.Second)
	assert.Len(mockBreaker.(*breaker).clientMap, 0)
}

func mockHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := Upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot upgrade: %v", err), http.StatusInternalServerError)
	}
	mockBreaker.Register(conn)
}
