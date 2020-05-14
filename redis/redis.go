package redis

import (
	"encoding/json"
	"time"

	"github.com/callicoder/go-commons/errors"
	"github.com/go-redis/redis/v7"
)

var (
	errInvalidOptions = errors.New("Invalid Options")
	errConnecting     = errors.New("Error connecting to Redis")
	ErrCacheMiss      = errors.New("Cache miss")
)

const (
	defaultDialTimeout  = 3000
	defaultReadTimeout  = 3000
	defaultWriteTimeout = 3000
)

type Config struct {
	Addrs        []string
	Pwd          string
	DB           int
	DialTimeout  int
	ReadTimeout  int
	WriteTimeout int
}

type RedisClient interface {
	Set(key string, value interface{}, ttl time.Duration) error
	SetStruct(key string, value interface{}, ttl time.Duration) error
	GetBytes(key string) ([]byte, error)
	GetString(key string) (string, error)
	GetStruct(key string, v interface{}) error
	Keys(pattern string) []string
	Del(key string) error

	HSet(key, field string, value interface{}) error
	HSetTTL(key, field string, value interface{}, t time.Duration)
	HGet(key, field string) (string, error)
	HGetAll(key string) map[string]string
	HDel(key, field string) error

	Close() error
	ClusterMode() bool
	Single() *redis.Client
	Cluster() *redis.ClusterClient
}

type redisClient struct {
	single      *redis.Client
	cluster     *redis.ClusterClient
	clusterMode bool
}

func NewRedisClient(c Config) (*redisClient, error) {
	if len(c.Addrs) == 0 {
		return nil, errors.Wrap(errInvalidOptions, "No Redis Addresses to connect...")
	}

	if c.DialTimeout == 0 {
		c.DialTimeout = defaultDialTimeout
	}

	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}

	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}

	r := &redisClient{}

	if len(c.Addrs) == 1 {
		r.single = redis.NewClient(
			&redis.Options{
				Addr:         c.Addrs[0],
				Password:     c.Pwd,
				DB:           c.DB,
				DialTimeout:  time.Duration(c.DialTimeout) * time.Millisecond,
				ReadTimeout:  time.Duration(c.ReadTimeout) * time.Millisecond,
				WriteTimeout: time.Duration(c.WriteTimeout) * time.Millisecond,
			})

		if err := r.single.Ping().Err(); err != nil {
			return nil, errConnecting
		}
		r.clusterMode = false
		return r, nil
	}

	r.cluster = redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:        c.Addrs,
			Password:     c.Pwd,
			DialTimeout:  time.Duration(c.DialTimeout) * time.Millisecond,
			ReadTimeout:  time.Duration(c.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(c.WriteTimeout) * time.Millisecond,
		})

	if err := r.cluster.Ping().Err(); err != nil {
		return nil, errConnecting
	}
	r.clusterMode = true

	return r, nil
}

func (r *redisClient) Set(key string, value interface{}, ttl time.Duration) error {
	if r.clusterMode {
		return r.cluster.Set(key, value, ttl).Err()
	}
	return r.single.Set(key, value, ttl).Err()
}

func (r *redisClient) SetStruct(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if r.clusterMode {
		r.cluster.Set(key, jsonData, ttl).Err()
	}

	return r.single.Set(key, jsonData, ttl).Err()
}

func (r *redisClient) GetBytes(key string) ([]byte, error) {
	var (
		value []byte
		err   error
	)

	if r.clusterMode {
		value, err = r.cluster.Get(key).Bytes()
	} else {
		value, err = r.single.Get(key).Bytes()
	}

	if err != nil {
		return nil, cacheError(err)
	}

	return value, nil
}

func (r *redisClient) GetString(key string) (string, error) {
	var (
		value string
		err   error
	)

	if r.clusterMode {
		value, err = r.cluster.Get(key).Result()
	} else {
		value, err = r.single.Get(key).Result()
	}

	if err != nil {
		return "", cacheError(err)
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

func (r *redisClient) Keys(pattern string) []string {
	if r.clusterMode {
		return r.cluster.Keys(pattern).Val()
	}
	return r.single.Keys(pattern).Val()
}

func (r *redisClient) Del(key string) error {
	if r.clusterMode {
		return r.cluster.Del(key).Err()
	}
	return r.single.Del(key).Err()
}

func (r *redisClient) HSet(key, field string, value interface{}) error {
	if r.clusterMode {
		return r.cluster.HSet(key, field, value).Err()
	}
	return r.single.HSet(key, field, value).Err()
}

func (r *redisClient) HSetTTL(key, field string, value interface{}, ttl time.Duration) error {
	if r.clusterMode {
		pipe := r.cluster.TxPipeline()
		pipe.HSet(key, field, value)
		pipe.Expire(key, ttl)

		_, err := pipe.Exec()
		return err
	}
	pipe := r.single.TxPipeline()
	pipe.HSet(key, field, value)
	pipe.Expire(key, ttl)

	_, err := pipe.Exec()
	return err
}

func (r *redisClient) HGet(key, field string) (string, error) {
	var (
		value string
		err   error
	)

	if r.clusterMode {
		value, err = r.cluster.HGet(key, field).Result()
	} else {
		value, err = r.single.HGet(key, field).Result()
	}

	if err != nil {
		return "", cacheError(err)
	}

	return value, nil
}

func (r *redisClient) HGetAll(k string) map[string]string {
	if r.clusterMode {
		return r.cluster.HGetAll(k).Val()
	}
	return r.single.HGetAll(k).Val()
}

func (r *redisClient) HDel(k, field string) error {
	if r.clusterMode {
		return r.cluster.HDel(k, field).Err()
	}
	return r.single.HDel(k, field).Err()
}

func (r *redisClient) ClusterMode() bool {
	return r.clusterMode
}

func (r *redisClient) Single() *redis.Client {
	return r.single
}

func (r *redisClient) Cluster() *redis.ClusterClient {
	return r.cluster
}

func (r *redisClient) Close() error {
	if r.clusterMode {
		return r.cluster.Close()
	}
	return r.single.Close()
}

func cacheError(err error) error {
	if err == redis.Nil {
		return ErrCacheMiss
	}
	return err
}
