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

	var args []interface{}
	for i := 0; i < 1000; i++ {
		args = append(args, fmt.Sprintf("a-%d-b", i), fmt.Sprintf("%d", i))
	}
	reply, err = m.Do("MSET", args...)
	assert.NoError(err)
	assert.Len(reply.([]interface{}), 3)
	for _, ok := range reply.([]interface{}) {
		assert.Equal("OK", ok)
	}

	reply, err = m.Do("MGET", args...)
	assert.NoError(err)
	var total []string
	for _, sep := range reply.([]interface{}) {
		result, err := redigo.Values(sep, err)
		assert.NoError(err)

		var uniq []string
		for i := 0; i < len(result); i += 2 {
			if result[i] == nil {
				continue
			}
			uniq = append(uniq, string(result[i].([]byte)))
		}
		total = append(total, uniq...)
	}
	assert.Equal(1000, len(total))

	var subargs []interface{}
	for i := 1000; i < 2000; i++ {
		subargs = append(subargs, fmt.Sprintf("a-%d-b", i), float64(i))
	}

	total = []string{}
	reply, err = m.Do("MGET", subargs...)
	assert.NoError(err)
	for _, sep := range reply.([]interface{}) {
		result, err := redigo.Values(sep, err)
		assert.NoError(err)

		var uniq []string
		for i := 0; i < len(result); i += 2 {
			if result[i] == nil {
				continue
			}
			uniq = append(uniq, string(result[i].([]byte)))
		}
		total = append(total, uniq...)
	}
	assert.Equal(0, len(total))

	total = []string{}
	newArgs := append(args, subargs...)
	reply, err = m.Do("MGET", newArgs...)
	assert.NoError(err)
	for _, sep := range reply.([]interface{}) {
		result, err := redigo.Values(sep, err)
		assert.NoError(err)
		var uniq []string
		for i := 0; i < len(result); i += 2 {
			if result[i] == nil {
				continue
			}
			uniq = append(uniq, string(result[i].([]byte)))
		}
		total = append(total, uniq...)
	}
	assert.Equal(1000, len(total))

	// Close
	err = m.Close()
	assert.NoError(err)
	assert.Len(group.Pools, 0)
	assert.Len(group.Ring.Members(), 0)
}
