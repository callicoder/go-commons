package server

type Config struct {
	ContextPath             string
	Port                    int
	ReadTimeout             int
	WriteTimeout            int
	GracefulShutdownTimeout int
	CORS                    CORSConfig
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedHeaders []string
	AllowedMethods []string
}
