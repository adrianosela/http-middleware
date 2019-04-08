package middleman

import (
	"context"
	"fmt"
	"net/http"
)

// BasicAuthenticator is a Middleware which runs a given function to
// authenticate requests with basic authentication in the appropriate header
type BasicAuthenticator struct {
	authenticate AuthenticateFunc
}

// AuthenticateFunc authenticates a user with basic credentials
type AuthenticateFunc func(uname, pw string) (string, error)

var (
	unameCtxKey = requestCtxKey("basic-auther-username")
)

// NewBasicAuthenticator is the constructor for a BasicAuthenticator Middleware
func NewBasicAuthenticator(af AuthenticateFunc) *BasicAuthenticator {
	return &BasicAuthenticator{authenticate: af}
}

// Wrap makes BasicAuthenticator implement the Middleware interface
func (ba *BasicAuthenticator) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uname, pw, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no basic credentials in request header"))
			return
		}
		uname, err := ba.authenticate(uname, pw)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("authentication failed. %s", err)))
			return
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), unameCtxKey, uname)))
	})
}

// GetUname returns the authenticated username
func GetUname(r *http.Request) string {
	return r.Context().Value(unameCtxKey).(string)
}
