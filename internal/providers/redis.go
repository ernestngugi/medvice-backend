package providers

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

type (
	Redis interface {
		Exists(key string) (bool, error)
		Get(key string) (interface{}, error)
		Set(key string, val interface{}) (interface{}, error)
	}

	RedisConfig struct {
		IdleTimeout time.Duration
		MaxActive   int
		MaxIdle     int
	}

	AppRedis struct {
		pool *redis.Pool
	}
)

func NewRedisProvider(
	config *RedisConfig,
) *AppRedis {
	return NewRedisWithURL(fmt.Sprintf(
		"redis://:%s@%s:%s",
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	), config)
}

func NewRedisWithURL(
	redisURL string,
	config *RedisConfig,
) *AppRedis {
	idleTimeout := 1 * time.Minute
	maxActive := 10
	maxIdle := 5

	if config != nil {
		if int64(config.IdleTimeout) != 0 {
			idleTimeout = config.IdleTimeout
		}

		if config.MaxActive != 0 {
			maxActive = config.MaxActive
		}

		if config.MaxIdle != 0 {
			maxIdle = config.MaxIdle
		}
	}

	redisPool := &redis.Pool{
		IdleTimeout: idleTimeout,
		MaxActive:   maxActive,
		MaxIdle:     maxIdle,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisURL)
		},
	}

	return &AppRedis{
		pool: redisPool,
	}
}

func (p *AppRedis) Exists(key string) (bool, error) {
	return redis.Bool(p.do("EXISTS", key))
}

func (p *AppRedis) Get(
	key string,
) (interface{}, error) {
	return p.do("GET", key)
}

func (p *AppRedis) Set(
	key string,
	val interface{},
) (interface{}, error) {
	return p.do("SET", key, val)
}

func (p *AppRedis) do(
	commandName string,
	args ...interface{},
) (interface{}, error) {
	conn := p.pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}
