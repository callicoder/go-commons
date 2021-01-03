package server

type Config struct {
	ContextPath               string
	Port                      int
	ReadTimeoutMs             int `mapstructure:"read_timeout_ms"`
	WriteTimeoutMs            int `mapstructure:"write_timeout_ms"`
	GracefulShutdownTimeoutMs int `mapstructure:"graceful_shutdown_timeout_ms"`
	CORS                      CORSConfig
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	MaxAge         int      `mapstructure:"maxage"`
}
