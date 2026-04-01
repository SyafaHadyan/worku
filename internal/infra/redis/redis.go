// Package redis connects and test the Redis configuration and provides interface
package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/redis/go-redis/v9"
)

type RedisItf interface {
	Set(key string, value string)
	Get(key string) (string, error)
	Delete(key string)
}

type Redis struct {
	Client     *redis.Client
	env        string
	expiration int
}

func New(env *env.Env) *Redis {
	url := fmt.Sprintf(
		"redis://%s:%s@%s:%d/%d",
		env.RedisUsername,
		env.RedisPassword,
		env.RedisAddress,
		env.RedisPort,
		env.RedisDatabase,
	)

	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Panic(err)
	}

	redis := redis.NewClient(opts)

	Redis := Redis{
		Client:     redis,
		expiration: env.RedisExpiration,
	}

	return &Redis
}

func Test(r *Redis) {
	ctx := context.Background()

	log.Println("testing redis connection")

	status := r.Client.Ping(ctx)

	if status.Err() != nil {
		log.Panic("failed to test redis instance")
	}

	log.Println("redis testing success")
}

func (r *Redis) Set(key string, value string) {
	ctx := context.Background()

	duration := time.Duration(r.expiration) * time.Minute

	r.Client.Set(ctx, key, value, duration)
}

func (r *Redis) Get(key string) (string, error) {
	value, err := r.Client.Get(context.Background(), key).Result()

	return value, err
}

func (r *Redis) Delete(key string) {
	_ = r.Client.Del(context.Background(), key)
}
