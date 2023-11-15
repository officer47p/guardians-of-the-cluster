package cache

import (
	"errors"
	"strconv"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() (*RedisCache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil || pong != "PONG" {
		return nil, errors.New("error while connecting to redis server")
	}

	return &RedisCache{client: client}, nil
}

func (r *RedisCache) FlushData() error {
	res, err := r.client.FlushAll().Result()
	if err != nil {
		return err
	}
	if res != "OK" {
		return errors.New("error when flushing the keys in redis")
	}

	return nil
}

func (r *RedisCache) GetKey(k string) (int64, error) {
	v, err := r.client.Get(k).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, ErrKeyNotFound
		}

		return 0, err
	}
	// v, ok := r.records[k]
	// if !ok {
	// 	return 0, ErrKeyNotFound
	// }

	vInt, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, errors.New("could not parse value to int")
	}

	return vInt, nil
}

func (r *RedisCache) SetKey(k string, v int64) error {
	err := r.client.Set(k, v, 0).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		return errors.New("error when setting the key")
	}

	return nil
}
