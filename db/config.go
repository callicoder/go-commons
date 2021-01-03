package db

import (
	"bytes"
	"strconv"
)

type Config struct {
	Driver             string
	Name               string
	Host               string
	Port               int
	Username           string
	Password           string
	Query              string
	MaxIdleConnections int `mapstructure:"max_idle_connections"`
	MaxOpenConnections int `mapstructure:"max_open_connections"`
}

func (cfg Config) URL() string {
	var buf bytes.Buffer

	buf.WriteString(cfg.Driver)
	buf.WriteString("://")

	// [username[:password]@]
	if len(cfg.Username) > 0 {
		buf.WriteString(cfg.Username)
		if len(cfg.Password) > 0 {
			buf.WriteByte(':')
			buf.WriteString(cfg.Password)
		}
		buf.WriteByte('@')
	}

	// [host[:port]]
	if len(cfg.Host) > 0 {
		buf.WriteString(cfg.Host)
		if cfg.Port > 0 {
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(cfg.Port))
		}
	}

	// /dbname
	buf.WriteByte('/')
	buf.WriteString(cfg.Name)

	// ?query=value
	if len(cfg.Query) > 0 {
		buf.WriteByte('?')
		buf.WriteString(cfg.Query)
	}

	return buf.String()
}
