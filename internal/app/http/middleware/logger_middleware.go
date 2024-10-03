package middleware

import (
	"net/http"
	"time"

	"github.com/luiidev/go/pkg/logger"
)

// Logger is a middleware handler that does request logging
type Logger struct {
	handler http.Handler
	l       *logger.Logger
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	l.l.Info("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

// NewLogger constructs a new Logger middleware handler
func NewLogger(handlerToWrap http.Handler, l *logger.Logger) *Logger {
	return &Logger{
		handler: handlerToWrap,
		l:       l,
	}
}
