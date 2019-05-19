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

func TestNewBreaker(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		inputs []Option
		wants  []int
	}{
		"empty": {inputs: []Option{}, wants: []int{512, 100}},
		"step1": {inputs: []Option{WithMaxReadLimit(100),
			WithMaxMessagePoolLength(10)}, wants: []int{100, 10}},
	}

	for _, test := range tests {
		bk, err := NewBreaker(test.inputs...)
		assert.NoError(err)
		bk.(*breaker).errorHandler(fmt.Errorf("[check] new breaker test"))
		assert.Equal(test.wants[0], int(bk.(*breaker).maxReadLimit))
		assert.Equal(test.wants[1], cap(bk.(*breaker).broadcast))
	}
}

func TestBreaker_Register(t *testing.T) {
	assert := assert.New(t)
	bk, err := NewBreaker(WithMaxReadLimit(1024), WithMaxMessagePoolLength(200))
	assert.NoError(err)

	srv := httptest.NewServer(http.HandlerFunc(testHandler))
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(err)
	err = bk.Register(conn)
	assert.NoError(err)
	time.Sleep(1 * time.Second)
	assert.Len(bk.(*breaker).clientMap, 1)
}

func TestBreaker_UnRegister(t *testing.T) {
	assert := assert.New(t)
	bk, err := NewBreaker(WithMaxReadLimit(1024), WithMaxMessagePoolLength(200))
	assert.NoError(err)

	srv := httptest.NewServer(http.HandlerFunc(testHandler))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(err)
	err = bk.Register(conn)
	assert.NoError(err)
	time.Sleep(1 * time.Second)
	assert.Len(bk.(*breaker).clientMap, 1)
	for k, _ := range bk.(*breaker).clientMap {
		err = bk.UnRegister(k)
		assert.NoError(err)
	}
	time.Sleep(1 * time.Second)
	assert.Len(bk.(*breaker).clientMap, 0)
}

func TestBreaker_BroadCast(t *testing.T) {
	assert := assert.New(t)

	bk, err := NewBreaker(WithMaxReadLimit(1024), WithMaxMessagePoolLength(200))
	assert.NoError(err)

	srv := httptest.NewServer(http.HandlerFunc(testHandler))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(err)
	err = bk.Register(conn)
	assert.NoError(err)
	time.Sleep(1 * time.Second)

	err = bk.BroadCast(&internalMessage{payload: []byte("hello world")})
	assert.NoError(err)
	time.Sleep(1 * time.Second)
}

func testHandler(w http.ResponseWriter, req *http.Request) {
	_, err := Upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot upgrade: %v", err), http.StatusInternalServerError)
	}
}
