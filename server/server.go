package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

type Server struct {
	httpServer *http.Server
	config     Config
}

func New(conf Config, handler http.Handler) *Server {
	allowedMethods := handlers.AllowedMethods(conf.CORS.AllowedMethods)
	allowedHeaders := handlers.AllowedHeaders(conf.CORS.AllowedHeaders)
	allowedOrigins := handlers.AllowedOrigins(conf.CORS.AllowedOrigins)
	maxAge := handlers.MaxAge(conf.CORS.MaxAge)
	handler = handlers.CORS(allowedMethods, allowedHeaders, allowedOrigins, maxAge)(handler)

	return &Server{
		httpServer: &http.Server{
			Handler:      handler,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			Addr:         fmt.Sprintf("0.0.0.0:%d", conf.Port),
		},
		config: conf,
	}
}

func (s *Server) Start() {
	go func() {
		log.Printf("Starting http server on port %v", s.config.Port)
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start http server %v", err)
			return
		}
	}()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.GracefulShutdownTimeout)*time.Millisecond)
	defer cancel()
	s.httpServer.Shutdown(ctx)
	log.Println("Shutting down http server.")
}
