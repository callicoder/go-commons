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
