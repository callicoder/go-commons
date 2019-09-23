package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigAddr(t *testing.T) {
	t.Run("It should return Redis Addr", func(t *testing.T) {
		redisConfig := Config{
			Host: "localhost",
			Port: 6379,
		}

		expectedURL := "localhost:6379"
		actualURL := redisConfig.Addr()
		assert.Equal(t, expectedURL, actualURL)
	})
}
