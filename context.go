package don

import (
	"context"
	"net/http"
)

type contextKey uint8

const (
	requestContextKey contextKey = iota
	responseContextKey
)

// Request returns the raw request.
func Request(ctx context.Context) *http.Request {
	return ctx.Value(requestContextKey).(*http.Request)
}

func WithRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestContextKey, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ResponseWriter returns the response writer.
func ResponseWriter(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseContextKey).(http.ResponseWriter)
}

func WithResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), responseContextKey, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
