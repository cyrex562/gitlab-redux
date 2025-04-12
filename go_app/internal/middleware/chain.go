package middleware

import "net/http"

// Middleware represents a middleware function
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain chains multiple middleware functions together
func Chain(middlewares ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(middlewares) - 1; i >= 0; i-- {
				last = middlewares[i](last)
			}
			last(w, r)
		}
	}
} 