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

// ResponseWriter returns the response writer.
func ResponseWriter(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseContextKey).(http.ResponseWriter)
}
