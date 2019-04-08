package middleman

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/teris-io/shortid"
)

const (
	loggerHeader = "[http]"
)

// Logger is a Middleware which logs all requests
type Logger struct {
	logDuration      bool
	logContentLength bool
	logStatus        bool
}

// LoggerConfig determines what the middleware will log
type LoggerConfig struct {
	LogDuration      bool
	LogContentLength bool
	LogStatus        bool
}

// NewLogger is the constructor for a Logger Middleware
func NewLogger(c LoggerConfig) *Logger {
	return &Logger{
		logDuration:      c.LogDuration,
		logContentLength: c.LogContentLength,
		logStatus:        c.LogStatus,
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// Wrap makes Logger implement the Middleware interface
func (l *Logger) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid, _ := shortid.Generate()
		start := time.Now().UnixNano()
		sw := statusWriter{ResponseWriter: w}

		logBefore := fmt.Sprintf("%s[%s] %s %s RECEIVED", loggerHeader, reqid, r.Method, r.URL.Path)
		logAfter := fmt.Sprintf("%s[%s] COMPLETED", loggerHeader, reqid)

		log.Print(logBefore)
		h.ServeHTTP(&sw, r)

		if l.logStatus {
			logAfter = fmt.Sprintf("%s with status %d %s", logAfter, sw.status, http.StatusText(sw.status))
		}
		if l.logDuration {
			logAfter = fmt.Sprintf("%s after %d ms", logAfter, (time.Now().UnixNano()-start)/1000000)
		}
		if l.logContentLength {
			logAfter = fmt.Sprintf("%s, Content-Length %d", logAfter, sw.length)
		}

		log.Print(logAfter)
	})
}
