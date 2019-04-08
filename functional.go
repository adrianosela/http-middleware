package middleman

import "net/http"

// Functional is a Middleware which runs given functions before and
// after every request
type Functional struct {
	before func()
	after  func()
}

// NewFunctional is the constructor for a Functional Middleware
func NewFunctional(before, after func()) *Functional {
	return &Functional{
		before: before,
		after:  after,
	}
}

// Wrap makes Functional implement the Middleware interface
func (f *Functional) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.before()
		h.ServeHTTP(w, r)
		f.after()
	})
}
