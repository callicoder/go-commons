package statsd

import (
	"testing"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/stretchr/testify/assert"
)

func TestReporter(t *testing.T) {
	t.Run("should configure client with given tags and namespace", func(t *testing.T) {
		var testConfig = Config{
			Host:      "localhost",
			Port:      6359,
			Tags:      []string{"error"},
			Namespace: "test",
			Enabled:   true,
		}

		reporter, err := NewClient(testConfig)

		assert.NoError(t, err)
		assert.Equal(t, []string{"error"}, reporter.Client.(*statsd.Client).Tags)
		assert.Equal(t, "test", reporter.Client.(*statsd.Client).Namespace)
	})

	t.Run("should configure create a NoOp client", func(t *testing.T) {
		var testConfig = Config{
			Host:      "localhost",
			Port:      6359,
			Tags:      []string{"error"},
			Namespace: "test",
			Enabled:   false,
		}

		reporter, err := NewClient(testConfig)

		assert.NoError(t, err)
		assert.Equal(t, &statsd.NoOpClient{}, reporter.Client.(*statsd.NoOpClient))
	})
}
