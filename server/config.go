package server

type Config struct {
	ContextPath                string
	Port                       int
	ReadTimeoutSec             int
	WriteTimeoutSec            int
	GracefulShutdownTimeoutSec int
	CORSAllowedOrigins         []string
}
