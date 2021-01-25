package redis

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	ctx = context.Background()
)

type RedisTestSuite struct {
	suite.Suite
	client        Client
	clusterClient Client
}

func (suite *RedisTestSuite) SetupSuite() {
	redisAddrs := os.Getenv("REDIS_ADDRS")
	if redisAddrs == "" {
		redisAddrs = "localhost:6379"
	}

	redisClusterAddrs := os.Getenv("REDIS_CLUSTER_ADDRS")
	if redisClusterAddrs == "" {
		log.Println("REDIS_CLUSTER_ADDRS not specified. Cluster functionality can't be tested")
		redisClusterAddrs = redisAddrs
	}

	cfg := Config{
		Addrs: []string{redisAddrs},
		Pwd:   "",
		DB:    0,
	}

	clusterCfg := Config{
		Addrs: strings.Split(redisClusterAddrs, ","),
		Pwd:   "",
		DB:    0,
	}

	client, err := NewClient(cfg)
	if err != nil {
		suite.T().Fatalf("Error creating redis client: %s", err)
	}

	clusterClient, err := NewClient(clusterCfg)
	if err != nil {
		suite.T().Fatalf("Error creating redis cluster client: %s", err)
	}

	suite.client = client
	suite.clusterClient = clusterClient
}

func (suite *RedisTestSuite) TearDownTest() {
	suite.client.FlushDB(ctx)
	suite.clusterClient.FlushDB(ctx)
}

func (suite *RedisTestSuite) TearDownSuite() {
	suite.client.Close()
	suite.clusterClient.Close()
}

func (suite *RedisTestSuite) TestSetAndGet() {
	for _, client := range []Client{suite.client, suite.clusterClient} {
		suite.Run("It should set and get a key in redis", func() {
			key := "test"

			err := client.Set(ctx, key, "Hello World", 0).Err()
			suite.NoError(err)

			value, err := client.Get(ctx, key).Result()
			suite.NoError(err)
			suite.Equal("Hello World", value)
		})
	}
}

func (suite *RedisTestSuite) TestHSetAndHGet() {
	for _, client := range []Client{suite.client, suite.clusterClient} {
		suite.Run("It should HSet and HGet a key in redis", func() {
			err := client.HSet(ctx, "user", "name", "Sachin").Err()
			suite.NoError(err)

			err = client.HSet(ctx, "user", "age", 25).Err()
			suite.NoError(err)

			value, err := client.HGet(ctx, "user", "name").Result()
			suite.NoError(err)
			suite.Equal("Sachin", value)

			value, err = client.HGet(ctx, "user", "age").Result()
			suite.NoError(err)
			suite.Equal("25", value)
		})
	}
}

func TestRedis(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
