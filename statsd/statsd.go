package statsd

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
)

type Config struct {
	Host      string
	Port      int
	Tags      []string
	Namespace string
	Enabled   bool
}

type Client interface {
	Increment(name string) error
	Decrement(name string) error
	IncrementWithTags(name string, tags ...string) error
	DecrementWithTags(name string, tags ...string) error
	IncrementBy(name string, value float64) error
	DecrementBy(name string, value float64) error
	Close() error
}

type Reporter struct {
	Client statsd.ClientInterface
}

func NewClient(c Config) (*Reporter, error) {
	if !c.Enabled {
		// Useful for testing
		client := &statsd.NoOpClient{}
		return &Reporter{Client: client}, nil
	}

	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	client, err := statsd.New(address)
	if err != nil {
		return nil, err
	}
	client.Namespace = c.Namespace
	client.Tags = c.Tags

	collector, err := plugins.InitializeStatsdCollector(&plugins.StatsdCollectorConfig{
		StatsdAddr: address,
		Prefix:     c.Namespace + ".hystrix",
	})
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("error initiating hystrix collector on statsD %+v", err)
	}
	metricCollector.Registry.Register(collector.NewStatsdCollector)

	return &Reporter{Client: client}, nil
}

func (r *Reporter) Increment(name string) error {
	return r.Client.Incr(name, nil, 1)
}

func (r *Reporter) Decrement(name string) error {
	return r.Client.Decr(name, nil, 1)
}

func (r *Reporter) IncrementWithTags(name string, tags ...string) error {
	return r.Client.Incr(name, tags, 1)
}

func (r *Reporter) DecrementWithTags(name string, tags ...string) error {
	return r.Client.Decr(name, tags, 1)
}

func (r *Reporter) IncrementBy(name string, value float64) error {
	return r.Client.Incr(name, nil, value)
}

func (r *Reporter) DecrementBy(name string, value float64) error {
	return r.Client.Decr(name, nil, value)
}

func (r *Reporter) Close() error {
	return r.Client.Close()
}
