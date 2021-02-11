package response

import (
	"encoding/json"
	"net/http"

	"github.com/callicoder/go-commons/logger"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

type httpError struct {
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.Errorf("Failed to write response: %v", err)
	}
}

func Error(w http.ResponseWriter, statusCode int, err error) {
	httpErr := httpError{
		Message: err.Error(),
	}
	JSON(w, statusCode, httpErr)
}
