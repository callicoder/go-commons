package redis

import (
	"context"
	"encoding/json"
	"time"

	"errors"

	"github.com/go-redis/redis/v8"
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

	SetStruct(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	ReadStruct(ctx context.Context, key string, v interface{}) error
	HSetEx(ctx context.Context, key, field string, value interface{}, t time.Duration) error
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

	ctx := context.Background()

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

		if err := r.Ping(ctx).Err(); err != nil {
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

	if err := r.Ping(ctx).Err(); err != nil {
		return nil, errConnecting
	}

	return r, nil
}

// Single client
func (r *client) SetStruct(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Set(ctx, key, jsonData, ttl).Err()
}

func (r *client) ReadStruct(ctx context.Context, key string, v interface{}) error {
	value, err := r.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(value, v)
}

func (r *client) HSetEx(ctx context.Context, key, field string, value interface{}, ttl time.Duration) error {
	pipe := r.TxPipeline()
	pipe.HSet(ctx, key, field, value)
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *client) ClusterMode() bool {
	return false
}

// Cluster client
func (r *clusterClient) SetStruct(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, jsonData, ttl).Err()
}

func (r *clusterClient) ReadStruct(ctx context.Context, key string, v interface{}) error {
	value, err := r.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(value, v)
}

func (r *clusterClient) HSetEx(ctx context.Context, key, field string, value interface{}, ttl time.Duration) error {
	pipe := r.TxPipeline()
	pipe.HSet(ctx, key, field, value)
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *clusterClient) ClusterMode() bool {
	return true
}
