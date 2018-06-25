package _interface

import "time"

type Consistence interface {
	Put(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	GetWithPrefix(key string) ([]string, error)
}
