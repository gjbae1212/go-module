package redis

import (
	"testing"

	"time"

	"fmt"

	"github.com/alicebob/miniredis"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	assert := assert.New(t)

	addrs := []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"}
	for _, addr := range addrs {
		s := miniredis.NewMiniRedis()
		if err := miniredis.NewMiniRedis().StartAddr(addr); err != nil {
			panic(err)
		}
		defer s.Close()
	}

	m, err := NewManager(addrs[:len(addrs)-1])
	assert.NoError(err)

	group := m.(*Group)
	assert.Equal(len(group.Pools), 2)

	// add
	err = group.addPool("")
	assert.Error(err)
	assert.True(errorx.IsOfType(err, emptyError))

	err = group.addPool("127.0.0.1:7000")
	assert.Error(err)

	err = group.addPool(addrs[2])
	assert.NoError(err)

	assert.Equal(len(group.Pools), 3)

	// delete
	err = group.deletePool("")
	assert.Error(err)
	assert.True(errorx.IsOfType(err, emptyError))

	err = group.deletePool("127.0.0.1:7002")
	assert.NoError(err)

	err = group.deletePool("127.0.0.1:7002")
	assert.Error(err)

	err = group.deletePool("empty")
	assert.Error(err)

	assert.Equal(len(group.Pools), 2)

	// get pool
	_, err = group.getPoolByKey("")
	assert.Error(err)
	assert.True(errorx.IsOfType(err, emptyError))

	pool, err := group.getPoolByKey("allan")
	assert.NoError(err)

	for i := 0; i < 100; i++ {
		tmp, err := group.getPoolByKey("allan")
		assert.NoError(err)
		assert.Equal(pool.Addr, tmp.Addr)
	}

	_, err = group.getPoolByAddr("")
	assert.Error(err)
	assert.True(errorx.IsOfType(err, emptyError))

	_, err = group.getPoolByAddr("111")
	assert.Error(err)
	assert.True(errorx.IsOfType(err, notFoundError))

	pool, err = group.getPoolByAddr("127.0.0.1:7001")
	assert.NoError(err)
	assert.Equal(pool.Addr, "127.0.0.1:7001")
}

func TestDo(t *testing.T) {
	assert := assert.New(t)

	addrs := []string{"127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"}
	for _, addr := range addrs {
		s := miniredis.NewMiniRedis()
		if err := miniredis.NewMiniRedis().StartAddr(addr); err != nil {
			panic(err)
		}
		defer s.Close()
	}

	m, err := NewManager(addrs)
	assert.NoError(err)

	group := m.(*Group)

	// single get or set
	reply, err := m.Do("GET", "empty")
	assert.NoError(err)
	assert.Nil(reply)

	reply, err = m.Do("SET", "allan", "1010")
	assert.NoError(err)
	assert.Equal("OK", reply)

	pool, err := group.getPoolByKey("allan")
	assert.NoError(err)

	stats, err := m.Stats()
	assert.NoError(err)
	assert.Equal(1, stats[pool.Addr].ActiveCount)
	assert.Equal(1, stats[pool.Addr].IdleCount)

	reply, err = m.Do("GET", "allan")
	assert.NoError(err)

	result1, err := redigo.Int(reply, err)
	assert.NoError(err)
	assert.Equal(1010, result1)

	reply, err = m.DoWithTimeout(time.Nanosecond, "SET", "allan1", 2020)
	assert.Error(err)

	reply, err = m.DoWithTimeout(time.Second*2, "SET", "allan1", 2020)
	assert.NoError(err)
	assert.Equal("OK", reply)

	// multiple get or set
	reply, err = m.Do("MSET", "allan", 1.0, "aaa")
	assert.Error(err)

	var setArgs []interface{}
	var getArgs []interface{}
	for i := 0; i < 1000; i++ {
		getArgs = append(getArgs, fmt.Sprintf("a-%d-b", i))
		setArgs = append(setArgs, fmt.Sprintf("a-%d-b", i), fmt.Sprintf("%d", i))
	}
	reply, err = m.Do("MSET", setArgs...)
	assert.NoError(err)
	assert.Len(reply.([]interface{}), 1000)
	for _, ok := range reply.([]interface{}) {
		assert.Equal("OK", ok)
	}

	reply, err = m.Do("MGET", getArgs...)
	assert.NoError(err)
	assert.Len(reply, 1000)
	for i := 0; i < 1000; i++ {
		v := string(reply.([]interface{})[i].([]byte))
		assert.Equal(fmt.Sprintf("%d", i), v)
	}

	var setArgs2 []interface{}
	var getArgs2 []interface{}
	for i := 1000; i < 2000; i++ {
		getArgs2 = append(getArgs2, fmt.Sprintf("a-%d-b", i))
		setArgs2 = append(setArgs2, fmt.Sprintf("a-%d-b", i), float64(i))
	}

	reply, err = m.Do("MGET", getArgs2...)
	assert.NoError(err)
	assert.Len(reply, 1000)
	for _, v := range reply.([]interface{}) {
		assert.Nil(v)
	}

	getArgs3 := append(getArgs2, getArgs...)
	reply, err = m.Do("MGET", getArgs3...)
	assert.NoError(err)
	assert.Len(reply, 2000)
	for i := 0; i < 2000; i++ {
		if i < 1000 {
			assert.Nil(reply.([]interface{})[i])
			continue
		}
		v := string(reply.([]interface{})[i].([]byte))
		assert.Equal(fmt.Sprintf("%d", i%1000), v)
	}

	// SET + MGET
	for i := 0; i < 10; i++ {
		reply, err = m.Do("SET", fmt.Sprintf("%d-", i), i)
		assert.NoError(err)
		assert.Equal("OK", reply)
	}
	reply, err = m.Do("MGET", "1-")
	assert.NoError(err)
	assert.Len(reply, 1)
	assert.Equal("1", string(reply.([]interface{})[0].([]byte)))

	reply, err = m.Do("MGET", "2-", "9-", "10-")
	assert.NoError(err)
	assert.Len(reply, 3)
	assert.Equal("2", string(reply.([]interface{})[0].([]byte)))
	assert.Equal("9", string(reply.([]interface{})[1].([]byte)))
	assert.Equal(nil, reply.([]interface{})[2])

	// Close
	err = m.Close()
	assert.NoError(err)
	assert.Len(group.Pools, 0)
	assert.Len(group.Ring.Members(), 0)
}
