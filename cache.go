package leikari

import (
	"errors"
	"sync"
)

type Cache interface {
	Set(key string, value interface{})
	Add(key string, value interface{}) error
	Replace(key string, value interface{}) error
	Get(key string) (interface{}, bool)
}

type cache struct {
	sync.RWMutex
	items map[string]interface{}
}

func NewCache() Cache {
	return &cache{
		items: make(map[string]interface{}),
	}
}

func (c *cache) Set(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.items[key] = value
}

func (c *cache) Add(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.items[key]; ok {
		return errors.New("item exists")
	}
	c.items[key] = value
	return nil
}

func (c *cache) Replace(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.items[key]; ok {
		c.items[key] = value
		return nil
	}
	return errors.New("item not exists")
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	if value, ok := c.items[key]; ok {
		return value, ok
	}
	return nil, false
}
