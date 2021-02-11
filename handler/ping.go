package handlers

import (
	"net/http"

	"github.com/callicoder/go-commons/logger"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("pong")); err != nil {
		logger.Errorf("Ping Error %v", err)
	}
}
