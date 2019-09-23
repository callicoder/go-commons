package redis

import (
	"fmt"
)

type Config struct {
	Network            string
	Host               string
	Port               int
	MaxIdleConnections int
	ReadTimeout        int
	ConnectionTimeout  int
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
