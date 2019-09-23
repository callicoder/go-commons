package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/callicoder/go-commons/errors"
	"github.com/go-redis/redis/v7"
)

type RedisClient interface {
	Set(key string, value interface{}, ttl time.Duration) error
	SetStruct(key string, value interface{}, ttl time.Duration) error
	GetBytes(key string) ([]byte, error)
	GetString(key string) (string, error)
	GetStruct(key string, v interface{}) error
	Close() error
	BaseClient() *redis.Client
}

type redisClient struct {
	baseClient *redis.Client
}

func NewRedisClient(conf Config) (*redisClient, error) {
	baseClient := redis.NewClient(&redis.Options{
		Addr:         conf.Addr(),
		DialTimeout:  time.Duration(conf.ConnectionTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(conf.ConnectionTimeout) * time.Millisecond,
		MinIdleConns: conf.MaxIdleConnections,
	})

	_, err := baseClient.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &redisClient{
		baseClient: baseClient,
	}, nil
}

func (r *redisClient) Set(key string, value interface{}, ttl time.Duration) error {
	return r.baseClient.Set(key, value, ttl).Err()
}

func (r *redisClient) SetStruct(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return r.baseClient.Set(key, jsonData, ttl).Err()
}

func (r *redisClient) GetBytes(key string) ([]byte, error) {
	value, err := r.baseClient.Get(key).Bytes()

	if err != nil {
		return nil, cacheError(key, err)
	}

	return value, nil
}

func (r *redisClient) GetString(key string) (string, error) {
	value, err := r.baseClient.Get(key).Result()

	if err != nil {
		return "", cacheError(key, err)
	}

	return value, nil
}

func (r *redisClient) GetStruct(key string, v interface{}) error {
	value, err := r.GetBytes(key)

	if err != nil {
		return err
	}

	return json.Unmarshal(value, v)
}

func (r *redisClient) Close() error {
	return r.baseClient.Close()
}

func cacheError(key string, err error) error {
	if err == redis.Nil {
		return errors.NewCacheMissError(err, fmt.Sprintf("key '%s' does not exist", key))
	}
	return err
}

func (r *redisClient) BaseClient() *redis.Client {
	return r.baseClient
}
