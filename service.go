package easyredis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// Service is a wrapper around go-redis for easy read and write
type Service struct {
	client *redis.Client
}

// New instantiates a new service.
func New(client *redis.Client) *Service {
	return &Service{
		client: client,
	}
}

// ConnectToRedis instantiates a new go-redis client.
func ConnectToRedis(host, port, db, password string) (*redis.Client, error) {
	const errMessage = "failed to initialize redis client"

	redisAddr := fmt.Sprintf("%s:%s", host, port)

	redisDB, err := strconv.Atoi(db)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       redisDB,
	})

	pingErr := redisClient.Ping(context.Background()).Err()
	if pingErr != nil {
		return nil, errors.Wrap(pingErr, errMessage)
	}

	return redisClient, nil
}

// Get returns a value for given key from redis.
//
// ctx should have a timeout
// key
func (s *Service) Get(ctx context.Context, key string) ([]byte, error) {
	const errMessage = "failed to get value from redis"

	result, err := s.client.Exists(ctx, key).Uint64()
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	existsInRedis := result == 1
	if !existsInRedis {
		return nil, errors.Wrap(ErrRedisKeyNotFound, errMessage)
	}

	val, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	return val, nil
}

// Set writes the given value to redis as json with the given TTL.
//
// ctx should have a timeout
//
// ttl is the time until the data should expire
//
// key is the key that is being used to store the data
//
// value can be any value
func (s *Service) Set(ctx context.Context, ttl time.Duration, key string, val interface{}) error {
	const errMessage = "failed to write value to redis"

	toCache, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, errMessage)
	}

	err = s.client.SetEX(ctx, key, toCache, ttl).Err()
	if err != nil {
		return errors.Wrap(err, errMessage)
	}

	return nil
}
