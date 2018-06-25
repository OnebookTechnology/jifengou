package consitence

import (
	"sync"
	"time"
)

type TempConsist struct {
	Store sync.Map
}

func (c *TempConsist) Put(key, value string, ttl time.Duration) error {
	c.Store.Store(key, value)
	return nil
}

func (c *TempConsist) Get(key string) (string, error) {
	v, ok := c.Store.Load(key)
	if !ok {
		return "", nil
	} else {
		return v.(string), nil
	}
}

func (c *TempConsist) Del(key string) error {
	c.Store.Delete(key)
	return nil
}

func (c *TempConsist) GetWithPrefix(key string) ([]string, error) {
	return nil, nil
}
