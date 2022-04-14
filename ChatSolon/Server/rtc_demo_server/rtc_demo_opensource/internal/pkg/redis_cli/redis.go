package redis_cli

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	DialTimeout     = 500 * time.Millisecond
	ReadTimeout     = 500 * time.Millisecond
	WriteTimeout    = 500 * time.Millisecond
	PoolTimeout     = 500 * time.Millisecond
	IdleTimeout     = 60 * time.Minute
	MinRetryBackoff = 8 * time.Millisecond
	MaxRetryBackoff = 128 * time.Millisecond
)

var Client *redis.Client

func NewRedis(addr, password string) {
	Client = redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        password,
		MaxRetries:      3,
		MinRetryBackoff: MinRetryBackoff,
		MaxRetryBackoff: MaxRetryBackoff,
		PoolSize:        100,
		DialTimeout:     DialTimeout,
		ReadTimeout:     ReadTimeout,
		WriteTimeout:    WriteTimeout,
		PoolTimeout:     PoolTimeout,
		IdleTimeout:     IdleTimeout,
	})

	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		panic("redis_cli init failed,error:" + err.Error())
	}

}
