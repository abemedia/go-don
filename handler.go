package don

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type StatusCoder interface {
	StatusCode() int
}

type Headerer interface {
	Header() http.Header
}

type Handle[T any] func(ctx context.Context, request T) (interface{}, error)

func H[T any](handle Handle[T]) http.Handler {
	return &handler[T]{handle: handle}
}

type handler[T any] struct {
	config *Config
	handle Handle[T]
}

func (h *handler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enc, err := getEncoder(r.Header.Get("Accept"), h.config.DefaultEncoding)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}

	res, err := h.handler(w, r)
	if err != nil {
		res = err
	}

	if h, ok := res.(Headerer); ok {
		headers := w.Header()
		for k, v := range h.Header() {
			headers[k] = v
		}
	}

	if sc, ok := res.(StatusCoder); ok {
		w.WriteHeader(sc.StatusCode())
	} else if res == nil {
		w.WriteHeader(http.StatusNoContent)
	}

	enc(w, res)
}

func (h *handler[T]) handler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var req T

	// Decode the header.
	if decodeHeader.Check(&req) {
		err := decodeHeader.Decode(&req, r.Header)
		if err != nil {
			return nil, err
		}
	}

	// Decode the URL query.
	if r.URL.RawQuery != "" && decodePath.Check(&req) {
		err := decodeQuery.Decode(&req, r.URL.Query())
		if err != nil {
			return nil, err
		}
	}

	// Decode the path params.
	p := httprouter.ParamsFromContext(r.Context())
	if len(p) != 0 && decodePath.Check(&req) {
		params := make(map[string][]string, len(p))
		for i := range p {
			params[p[i].Key] = []string{p[i].Value}
		}

		err := decodePath.Decode(&req, params)
		if err != nil {
			return nil, err
		}
	}

	// Decode the body.
	if r.ContentLength > 0 {
		dec, err := getDecoder(r.Header.Get("Content-Type"), h.config.DefaultEncoding)
		if err != nil {
			return nil, err
		}

		if err := dec(r, &req); err != nil {
			return nil, err
		}
	}

	ctx := context.WithValue(r.Context(), requestContextKey, r)
	ctx = context.WithValue(ctx, responseContextKey, w)

	return h.handle(ctx, req)
}

func (h *handler[T]) setConfig(r *Config) {
	h.config = r
}
