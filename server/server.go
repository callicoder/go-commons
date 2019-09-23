package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/callicoder/go-commons/logger"
)

type Server struct {
	httpServer *http.Server
	config     Config
}

func NewServer(conf Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Handler:      handler,
			ReadTimeout:  time.Duration(conf.ReadTimeoutSec) * time.Second,
			WriteTimeout: time.Duration(conf.WriteTimeoutSec) * time.Second,
			Addr:         fmt.Sprintf("0.0.0.0:%d", conf.Port),
		},
		config: conf,
	}
}

func (s *Server) Start() {
	go func() {
		logger.Infof("Starting http server on port %v", s.config.Port)
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start http server %v", err)
			return
		}
	}()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.GracefulShutdownTimeoutSec)*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx)
	logger.Info("Shutting down http server.")
}
