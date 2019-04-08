package httpmiddleware

import (
	"log"
	"net/http"

	"github.com/teris-io/shortid"
)

const (
	loggerHeader = "[http]"
)

// Logger Middleware logs all http requests served
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid, _ := shortid.Generate()
		log.Printf("%s %s %s >> RECEIVED %s", loggerHeader, r.Method, r.URL.Path, reqid)
		h.ServeHTTP(w, r)
		log.Printf("%s %s %s << DONE %s with status %s", loggerHeader, r.Method, r.URL.Path, reqid, r.Response.Status)
	})
}
