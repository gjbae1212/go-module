package redis

import (
	"time"

	"github.com/gjbae1212/consistent"
	redigo "github.com/gomodule/redigo/redis"
)

type Manager interface {
	Do(command string, args ...interface{}) (interface{}, error)
	DoWithTimeout(timeout time.Duration, command string, args ...interface{}) (interface{}, error)
	Close() error
	Stats() (map[string]redigo.PoolStats, error)
}

func NewManager(addrs []string) (Manager, error) {
	if len(addrs) == 0 {
		return nil, emptyError.New("NewManager")
	}
	g := &Group{
		Opt:   defaultOption,
		Ring:  consistent.New(),
		Pools: map[string]*Pool{},
	}
	for _, addr := range addrs {
		if err := g.addPool(addr); err != nil {
			return nil, err
		}
	}
	return Manager(g), nil
}
