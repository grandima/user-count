package storage

import (
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

var Error404 = errors.New("key not found")

type RedisStore struct {
	Client redis.Conn
}

func NewRedisCache(address string) *RedisStore {
	var cache = &RedisStore{}
	conn, err := redis.DialURL(address)
	if err != nil {
		log.Fatal(err)
	}
	cache.Client = conn
	return cache
}

func (r *RedisStore) Get(key string) (int, error) {

	val, err := redis.Int(r.Client.Do("GET", key))
	return val, err
}

func (r *RedisStore) Set(key string, value string) error {
	_, err := r.Client.Do("SETEX", key, "120", value)
	return errors.Wrap(err, "error: settin  with SETEX")
}

func (r *RedisStore) SetIfNotExistsToZero(key string) error {
	_, err := r.Client.Do("SETNX", key, 0)
	return errors.Wrap(err, "error: settin  with SETNXtoZero")
}

func (r *RedisStore) Increment(key string) (int, error) {
	res, err := redis.Int(r.Client.Do("INCR", key))
	return res, errors.Wrap(err, "error: settin  with SETNXtoZero")
}

func (r *RedisStore) Exists(key string) (int, error) {
	e, err := redis.Int(r.Client.Do("EXISTS", key))
	if err != nil {
		return e, errors.Wrap(err, "error while checkin if token exists")
	}

	return e, nil
}
