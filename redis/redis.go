package redis

import (
	"encoding/json"
	"time"

	"errors"

	"github.com/go-redis/redis/v7"
)

var (
	errInvalidOptions = errors.New("Invalid Options")
	errConnecting     = errors.New("Error connecting to Redis")
	ErrCacheMiss      = redis.Nil
)

type Config struct {
	Addrs        []string
	Pwd          string
	DB           int
	DialTimeout  int
	ReadTimeout  int
	WriteTimeout int
	PoolSize     int
}

type Client interface {
	redis.Cmdable

	SetStruct(key string, value interface{}, ttl time.Duration) error
	ReadStruct(key string, v interface{}) error
	HSetEx(key, field string, value interface{}, t time.Duration) error
	Close() error
	ClusterMode() bool
}

type client struct {
	*redis.Client
}

type clusterClient struct {
	*redis.ClusterClient
}

func NewClient(c Config) (Client, error) {
	if len(c.Addrs) == 0 {
		return nil, errInvalidOptions
	}

	if len(c.Addrs) == 1 {
		r := &client{}
		r.Client = redis.NewClient(
			&redis.Options{
				Addr:         c.Addrs[0],
				Password:     c.Pwd,
				DB:           c.DB,
				DialTimeout:  time.Duration(c.DialTimeout) * time.Millisecond,
				ReadTimeout:  time.Duration(c.ReadTimeout) * time.Millisecond,
				WriteTimeout: time.Duration(c.WriteTimeout) * time.Millisecond,
				PoolSize:     c.PoolSize,
			})

		if err := r.Ping().Err(); err != nil {
			return nil, errConnecting
		}
		return r, nil
	}

	r := &clusterClient{}
	r.ClusterClient = redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:        c.Addrs,
			Password:     c.Pwd,
			DialTimeout:  time.Duration(c.DialTimeout) * time.Millisecond,
			ReadTimeout:  time.Duration(c.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(c.WriteTimeout) * time.Millisecond,
			PoolSize:     c.PoolSize,
		})

	if err := r.Ping().Err(); err != nil {
		return nil, errConnecting
	}

	return r, nil
}

// Single client
func (r *client) SetStruct(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Set(key, jsonData, ttl).Err()
}

func (r *client) ReadStruct(key string, v interface{}) error {
	value, err := r.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(value, v)
}

func (r *client) HSetEx(key, field string, value interface{}, ttl time.Duration) error {
	pipe := r.TxPipeline()
	pipe.HSet(key, field, value)
	pipe.Expire(key, ttl)

	_, err := pipe.Exec()
	return err
}

func (r *client) ClusterMode() bool {
	return false
}

// Cluster client
func (r *clusterClient) SetStruct(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(key, jsonData, ttl).Err()
}

func (r *clusterClient) ReadStruct(key string, v interface{}) error {
	value, err := r.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(value, v)
}

func (r *clusterClient) HSetEx(key, field string, value interface{}, ttl time.Duration) error {
	pipe := r.TxPipeline()
	pipe.HSet(key, field, value)
	pipe.Expire(key, ttl)

	_, err := pipe.Exec()
	return err
}

func (r *clusterClient) ClusterMode() bool {
	return true
}
