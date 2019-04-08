package middleman

import "net/http"

// Middleware wraps before/after operations around an http handler
type Middleware interface {
	Wrap(h http.Handler) http.Handler
}

// Middleman is a collection of Middleware to be applied to a handler
type Middleman struct {
	middleware []Middleware
}

// NewMiddleman returns a Middleman with all the given Middleware
func NewMiddleman(middleware ...Middleware) *Middleman {
	return &Middleman{middleware: middleware}
}

// Wrap wraps an http handler with all Middleware in the Middleman
// Note that Middleman can be nested, as Middleman itself implements
// the Middleware interface
func (m *Middleman) Wrap(handler http.Handler) http.Handler {
	h := handler
	// wrap in reverse order
	for i := len(m.middleware) - 1; i >= 0; i-- {
		h = m.middleware[i].Wrap(h)
	}
	return h
}
