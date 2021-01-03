package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

type Server struct {
	server *http.Server
	config Config
}

func New(conf Config, handler http.Handler) *Server {
	allowedMethods := handlers.AllowedMethods(conf.CORS.AllowedMethods)
	allowedHeaders := handlers.AllowedHeaders(conf.CORS.AllowedHeaders)
	allowedOrigins := handlers.AllowedOrigins(conf.CORS.AllowedOrigins)
	maxAge := handlers.MaxAge(conf.CORS.MaxAge)
	handler = handlers.CORS(allowedMethods, allowedHeaders, allowedOrigins, maxAge)(handler)

	return &Server{
		server: &http.Server{
			Handler:      handler,
			ReadTimeout:  time.Duration(conf.ReadTimeoutMs) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeoutMs) * time.Millisecond,
			Addr:         fmt.Sprintf("0.0.0.0:%d", conf.Port),
		},
		config: conf,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("%w :: Failed to start listener on %s", err, s.server.Addr)
	}

	go func() {
		log.Printf("Starting http server on port %v", s.config.Port)
		err := s.server.Serve(lis)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start http server %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.GracefulShutdownTimeoutMs)*time.Millisecond)
	defer cancel()
	log.Printf("Shutting down http server")
	return s.server.Shutdown(ctx)
}
