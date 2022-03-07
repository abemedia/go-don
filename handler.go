package don

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// StatusCoder allows you to customise the HTTP response code.
type StatusCoder interface {
	StatusCode() int
}

// Headerer allows you to customise the HTTP headers.
type Headerer interface {
	Header() http.Header
}

// H wraps your handler function with the Go generics magic.
func H[T any, O any](handle Handle[T, O]) http.Handler {
	return &handler[T, O]{handle: handle}
}

// Handle is the type for your handlers.
type Handle[T any, O any] func(ctx context.Context, request *T) (O, error)

type handler[T any, O any] struct {
	config *Config
	handle Handle[T, O]
}

func (h *handler[T, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enc, err := getEncoder(r.Header.Get("Accept"), h.config.DefaultEncoding)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}

	req := new(T)

	// Decode the header.
	if decodeHeader.Check(req) {
		err := decodeHeader.Decode(req, r.Header)
		if err != nil {
			enc(w, err)
			return
		}
	}

	// Decode the URL query.
	if r.URL.RawQuery != "" && decodePath.Check(req) {
		err := decodeQuery.Decode(req, r.URL.Query())
		if err != nil {
			enc(w, err)
			return
		}
	}

	// Decode the path params.
	p := httprouter.ParamsFromContext(r.Context())
	if len(p) != 0 && decodePath.Check(req) {
		params := make(map[string][]string, len(p))
		for i := range p {
			params[p[i].Key] = []string{p[i].Value}
		}

		err := decodePath.Decode(req, params)
		if err != nil {
			enc(w, err)
			return
		}
	}

	// Decode the body.
	if r.ContentLength > 0 {
		dec, err := getDecoder(r.Header.Get("Content-Type"), h.config.DefaultEncoding)
		if err != nil {
			enc(w, err)
			return
		}

		if err := dec(r, req); err != nil {
			enc(w, err)
			return
		}
	}

	ctx := context.WithValue(r.Context(), requestContextKey, r)
	ctx = context.WithValue(ctx, responseContextKey, w)
	res, err := h.handle(ctx, req)
	if err != nil {
		enc(w, res)
		return
	}

	resAny := any(res)

	if h, ok := resAny.(Headerer); ok {
		headers := w.Header()
		for k, v := range h.Header() {
			headers[k] = v
		}
	}

	if sc, ok := resAny.(StatusCoder); ok {
		w.WriteHeader(sc.StatusCode())
	} else if resAny == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	enc(w, res)
}

func (h *handler[T, O]) setConfig(r *Config) {
	h.config = r
}
