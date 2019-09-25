package server

type Config struct {
	ContextPath                string
	Port                       int
	ReadTimeoutSec             int
	WriteTimeoutSec            int
	GracefulShutdownTimeoutSec int
	CORS                       CORSConfig
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedHeaders []string
	AllowedMethods []string
}
