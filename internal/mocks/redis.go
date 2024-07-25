package mocks

import (
	"github.com/gomodule/redigo/redis"
)

type (
	payload struct {
		Value interface{}
	}

	MockRedis struct {
		store map[string]*payload
	}
)

func NewMockRedisProvider() *MockRedis {
	return &MockRedis{
		store: make(map[string]*payload),
	}
}

func (p *MockRedis) Exists(key string) (bool, error) {
	_, err := p.Get(key)
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	return err == nil, nil
}

func (p *MockRedis) Get(key string) (interface{}, error) {
	payload, ok := p.store[key]
	if !ok {
		return nil, redis.ErrNil
	}

	return payload.Value, nil
}

func (p *MockRedis) Set(key string, val interface{}) (interface{}, error) {
	newPayload := &payload{
		Value: val,
	}

	p.store[key] = newPayload

	return val, nil
}

func (p *MockRedis) Del(key string) error {
	delete(p.store, key)
	return nil
}
