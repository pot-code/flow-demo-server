package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type logEntry struct {
	r *http.Request
}

func (l *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	log.Debug().
		Str("method", l.r.Method).
		Str("uri", l.r.RequestURI).
		Int("status", status).
		Dur("elapsed", elapsed).
		Msg("")
}

func (l *logEntry) Panic(v interface{}, stack []byte) {}

type logFormatter func(r *http.Request) *logEntry

func (f logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return f(r)
}

func defaultLogFormatter(r *http.Request) *logEntry {
	return &logEntry{r}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return middleware.RequestLogger(logFormatter(defaultLogFormatter))(next)
}
