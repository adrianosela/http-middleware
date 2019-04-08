package middleman

import (
	"fmt"
	"net/http"
)

// NotFoundRedirector is a Middleware which redirects requests which
// your handler responds to with a 404, this is particularly useful when
// hosting a web application and you want to redirect 404s to a custom landing page
type NotFoundRedirector struct {
	redirectPath string
}

// NewNotFoundRedirector is the constructor for a NotFoundRedirector Middleware
func NewNotFoundRedirector(redirectPath string) *NotFoundRedirector {
	return &NotFoundRedirector{redirectPath: redirectPath}
}

type nfRespWriter struct {
	http.ResponseWriter
	status int
}

func (w *nfRespWriter) WriteHeader(status int) {
	w.status = status
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *nfRespWriter) Write(b []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(b)
	}
	return len(b), nil // lie about having written
}

// Wrap makes NotFoundRedirector implement the Middleware interface
func (nfr *NotFoundRedirector) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &nfRespWriter{ResponseWriter: w}
		h.ServeHTTP(rw, r)
		if rw.status == http.StatusNotFound {
			http.Redirect(w, r, fmt.Sprintf("%s%s", r.URL.Host, nfr.redirectPath), http.StatusFound)
		}
	})
}
