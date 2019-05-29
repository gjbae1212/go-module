package redis

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gjbae1212/consistent"
	redigo "github.com/gomodule/redigo/redis"
)

var (
	defaultOption = &Option{
		Idle:            5,
		Active:          50,
		IdleTimeout:     300 * time.Second,
		ConnTimeout:     100 * time.Millisecond,
		ReadTimeout:     300 * time.Millisecond,
		WriteTimeout:    2000 * time.Millisecond,
		MaxConnLifetime: 1 * time.Hour,
		Wait:            false,
	}
)

type (
	Group struct {
		sync.RWMutex
		Opt   *Option
		Ring  *consistent.Consistent
		Pools map[string]*Pool
	}

	Pool struct {
		*redigo.Pool
		Addr string
	}

	Option struct {
		Idle            int           // max idle connection count
		Active          int           // max active connection count
		IdleTimeout     time.Duration // connection will remain for this duration
		ConnTimeout     time.Duration // tcp handshake timeout
		ReadTimeout     time.Duration // request read timeout
		WriteTimeout    time.Duration // read and response timeout
		MaxConnLifetime time.Duration // per connection max alive duration
		Wait            bool          // if pool is overflow at max active limit, whether wait
	}

	multiReply struct {
		args     []interface{}
		response interface{}
		err      error
	}
)

func (c *Group) Do(command string, args ...interface{}) (interface{}, error) {
	return c.do(0, command, args...)
}

func (c *Group) DoWithTimeout(timeout time.Duration, command string, args ...interface{}) (interface{}, error) {
	return c.do(timeout, command, args...)
}

func (c *Group) Close() error {
	c.Lock()
	defer c.Unlock()

	for _, addr := range c.Ring.Members() {
		pool := c.Pools[addr]
		c.Ring.Remove(addr)
		delete(c.Pools, addr)
		if err := pool.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Group) Stats() (map[string]redigo.PoolStats, error) {
	stats := map[string]redigo.PoolStats{}
	for _, addr := range c.Ring.Members() {
		pool := c.Pools[addr]
		stats[addr] = pool.Stats()
	}
	return stats, nil
}

func (c *Group) addPool(addr string) error {
	if addr == "" {
		return emptyError.New("addPool")
	}
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Pools[addr]; ok {
		return alreadyError.New("addPool")
	}

	c.Ring.Add(addr)
	c.Pools[addr] = newPool(addr, defaultOption)
	return nil
}

func (c *Group) deletePool(addr string) error {
	if addr == "" {
		return emptyError.New("deletePool")
	}

	c.Lock()
	defer c.Unlock()

	pool, ok := c.Pools[addr]
	if !ok {
		return notFoundError.New("deletePool")
	}

	c.Ring.Remove(addr)
	delete(c.Pools, addr)
	return pool.Close()
}

func (c *Group) getPoolByKey(key string) (*Pool, error) {
	if key == "" {
		return nil, emptyError.New("getPoolByKey")
	}

	c.RLock()
	defer c.RUnlock()

	addr, err := c.Ring.Get(key)
	if err != nil {
		return nil, err
	}

	return c.Pools[addr], nil
}

func (c *Group) getPoolByAddr(addr string) (*Pool, error) {
	if addr == "" {
		return nil, emptyError.New("getPoolByAddr")
	}

	c.RLock()
	defer c.RUnlock()

	if pool, ok := c.Pools[addr]; ok {
		return pool, nil
	}
	return nil, notFoundError.New("getPoolByAddr")
}

func (c *Group) do(timeout time.Duration, command string, args ...interface{}) (interface{}, error) {
	if command == "" || len(args) < 1 {
		return nil, emptyError.New("Do")
	}

	if strings.HasPrefix(command, "MSET") || strings.HasPrefix(command, "MGET") {
		return c.multi(timeout, command, args...)
	}

	key, err := keyString(args[0])
	if err != nil {
		return nil, err
	}

	pool, err := c.getPoolByKey(key)
	if err != nil {
		return nil, err
	}
	conn := pool.Get()
	defer conn.Close()

	if timeout == 0 {
		return conn.Do(command, args...)
	} else {
		return redigo.DoWithTimeout(conn, timeout, command, args...)
	}
}

func (c *Group) asyncDo(timeout time.Duration, ch chan<- *multiReply, pool *Pool, command string, args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			ch <- &multiReply{response: nil, err: err.(error)}
		}
	}()
	conn := pool.Get()
	defer conn.Close()

	var result interface{}
	var err error
	if timeout == 0 {
		result, err = conn.Do(command, args...)
	} else {
		result, err = redigo.DoWithTimeout(conn, timeout, command, args...)
	}

	ch <- &multiReply{args: args, response: result, err: err}
}

func (c *Group) multi(timeout time.Duration, command string, args ...interface{}) ([]interface{}, error) {
	if command == "" || len(args) == 0 {
		return nil, emptyError.New("multi")
	}

	if !strings.HasPrefix(command, "MSET") && !strings.HasPrefix(command, "MGET") {
		return nil, invalidError.New("multi")
	}

	if strings.HasPrefix(command, "MSET") && len(args)%2 != 0 {
		return nil, invalidError.New("multi")
	}

	group := map[*Pool][]interface{}{}
	var allKeys []interface{}
	switch {
	case strings.HasPrefix(command, "MSET"):
		for i := 0; i < len(args); i += 2 {
			key, err := keyString(args[i])
			if err != nil {
				return nil, err
			}
			pool, err := c.getPoolByKey(key)
			if err != nil {
				return nil, err
			}
			allKeys = append(allKeys, args[i])
			group[pool] = append(group[pool], args[i], args[i+1])
		}
	case strings.HasPrefix(command, "MGET"):
		for i := 0; i < len(args); i += 1 {
			key, err := keyString(args[i])
			if err != nil {
				return nil, err
			}
			pool, err := c.getPoolByKey(key)
			if err != nil {
				return nil, err
			}
			allKeys = append(allKeys, args[i])
			group[pool] = append(group[pool], args[i])
		}
	}

	var chs []chan *multiReply
	for pool, args := range group {
		ch := make(chan *multiReply)
		chs = append(chs, ch)
		go c.asyncDo(timeout, ch, pool, command, args...)
	}

	var reply []interface{}
	collectionMap := make(map[interface{}]interface{})
	for _, ch := range chs {
		callback := <-ch
		switch {
		case strings.HasPrefix(command, "MSET"):
			if callback.err != nil {
				return nil, callback.err
			}
			for i := 0; i < len(callback.args); i += 2 {
				collectionMap[callback.args[i]] = callback.response
			}
		case strings.HasPrefix(command, "MGET"):
			if _, err := redigo.Values(callback.response, callback.err); err != nil {
				return nil, err
			}
			for i := 0; i < len(callback.args); i += 1 {
				collectionMap[callback.args[i]] = callback.response.([]interface{})[i]
			}
		}
	}
	for _, key := range allKeys {
		reply = append(reply, collectionMap[key])
	}
	return reply, nil
}

func keyString(key interface{}) (string, error) {
	switch key := key.(type) {
	case int:
		return strconv.FormatInt(int64(key), 10), nil
	case int32:
		return strconv.FormatInt(int64(key), 10), nil
	case int64:
		return strconv.FormatInt(int64(key), 10), nil
	case float32:
		return strconv.FormatFloat(float64(key), 'G', -1, 64), nil
	case float64:
		return strconv.FormatFloat(key, 'G', -1, 64), nil
	case []byte:
		return string(key), nil
	case string:
		return key, nil
	default:
		return "", invalidError.New("keyString")
	}
}

func newPool(addr string, opt *Option) *Pool {
	return &Pool{
		Pool: &redigo.Pool{
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", addr,
					redigo.DialConnectTimeout(opt.ConnTimeout),
					redigo.DialReadTimeout(opt.ReadTimeout),
					redigo.DialWriteTimeout(opt.WriteTimeout),
				)
			},
			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
			MaxIdle:         opt.Idle,
			MaxActive:       opt.Active,
			IdleTimeout:     opt.IdleTimeout,
			Wait:            opt.Wait,
			MaxConnLifetime: opt.MaxConnLifetime,
		},
		Addr: addr,
	}
}
