package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"errors"
)

// Context key type — unexported to prevent collisions
type contextKey string

const RequestIDKey contextKey = "requestID"

// GetRequestID pulls request ID from context — usable by anyone
func GetRequestID(ctx context.Context) string {
    id, _ := ctx.Value(RequestIDKey).(string)
    return id
}

// AppError — structured error with HTTP code
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return e.Message
}

// ErrorResponse — JSON error shape
type ErrorResponse struct {
    Error string `json:"error"`
}

// AppHandler — handler that returns error
type AppHandler func(http.ResponseWriter, *http.Request) error

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    err := fn(w, r)
    if err == nil {
        return
    }

    var appErr *AppError
    if errors.As(err, &appErr) {
        log.Printf("app error %d: %v", appErr.Code, appErr.Err)
        WriteJSON(w, appErr.Code, ErrorResponse{Error: appErr.Message})
    } else {
        log.Printf("unexpected error: %v", err)
        WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
    }
}

// WriteJSON — shared response helper
func WriteJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}