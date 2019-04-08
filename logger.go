package middleman

import (
	"log"
	"net/http"
	"time"

	"github.com/teris-io/shortid"
)

const (
	loggerHeader = "[http]"
)

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

// Logger Middleware logs all http requests served
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid, _ := shortid.Generate()
		start := time.Now().UnixNano()
		sw := statusWriter{ResponseWriter: w}

		log.Printf("%s[%s] %s %s",
			loggerHeader,
			reqid,
			r.Method,
			r.URL.Path,
		)

		h.ServeHTTP(&sw, r)

		log.Printf("%s[%s] %d %s, took %d ms",
			loggerHeader,
			reqid,
			sw.status,
			http.StatusText(sw.status),
			(time.Now().UnixNano()-start)/1000000,
		)
	})
}
