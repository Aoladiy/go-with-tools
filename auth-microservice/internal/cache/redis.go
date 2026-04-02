package cache

import (
	"context"
	"log"
	"time"

	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func New(c config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.RdbAddr,
		Username:     c.RdbUsr,
		Password:     c.RdbPsw,
		DB:           c.RdbId,
		MaxRetries:   c.RdbMaxRetries,
		ReadTimeout:  c.RdbReadTimeout,
		WriteTimeout: c.RdbWriteTimeout,
		MinIdleConns: c.RdbMinIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("cannot connect to redis: %v", err)
	}

	return rdb
}
