package middleware

import (
	"sync"
	"time"
)

var FlowCounterHandler *FlowCounter

type FlowCounter struct {
	RedisFlowCountMap   map[string]*RedisFlowCountService
	RedisFlowCountSlice []*RedisFlowCountService
	Locker              sync.RWMutex
}

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap:   map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Locker:              sync.RWMutex{},
	}
}

func init() {
	FlowCounterHandler = NewFlowCounter()
}

func (c *FlowCounter) GetFlowCounter(serverName string) (*RedisFlowCountService, error) {
	for _, it := range c.RedisFlowCountSlice {
		if it.AppID == serverName {
			return it, nil
		}
	}
	newCounter := NewRedisFlowCountService(serverName, 1*time.Second)
	c.RedisFlowCountSlice = append(c.RedisFlowCountSlice, newCounter)
	c.Locker.Lock()
	defer c.Locker.Unlock()
	c.RedisFlowCountMap[serverName] = newCounter
	return newCounter, nil
}
