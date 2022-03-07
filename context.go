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

func Request(ctx context.Context) *http.Request {
	return ctx.Value(requestContextKey).(*http.Request)
}

func ResponseWriter(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseContextKey).(http.ResponseWriter)
}
