package don

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/abemedia/go-don/decoder"
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

type Handler interface {
	http.Handler
	handle(http.ResponseWriter, *http.Request, httprouter.Params)
}

// H wraps your handler function with the Go generics magic.
func H[T any, O any](handle Handle[T, O]) Handler {
	h := &handler[T, O]{handler: handle}

	var t T

	if hasTag(t, headerTag) {
		dec, err := decoder.NewCachedDecoder(t, headerTag)
		if err == nil {
			h.decodeHeader = dec
		}
	}

	if hasTag(t, queryTag) {
		dec, err := decoder.NewMapDecoder(t, queryTag)
		if err == nil {
			h.decodeQuery = dec
		}
	}

	if hasTag(t, pathTag) {
		dec, err := decoder.NewParamsDecoder(t, pathTag)
		if err == nil {
			h.decodePath = dec
		}
	}

	return h
}

// Handle is the type for your handlers.
type Handle[T any, O any] func(ctx context.Context, request T) (O, error)

type handler[T any, O any] struct {
	config       *Config
	handler      Handle[T, O]
	decodeHeader *decoder.CachedDecoder
	decodePath   *decoder.ParamsDecoder
	decodeQuery  *decoder.MapDecoder
	isNil        func(v any) bool
}

func (h *handler[T, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handle(w, r, nil)
}

//nolint:gocognit,cyclop
func (h *handler[T, O]) handle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	contentType := h.getEncoding(r.Header, "Accept")

	enc, err := getEncoder(contentType)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}

	req := new(T)

	var (
		res any
		e   *HTTPError
	)

	// Decode the header.
	if h.decodeHeader != nil {
		err = h.decodeHeader.Decode(r.Header, req)
		if err != nil {
			e = Error(err, http.StatusBadRequest)
			goto Encode
		}
	}

	// Decode the URL query.
	if h.decodeQuery != nil && r.URL.RawQuery != "" {
		err := h.decodeQuery.Decode(r.URL.Query(), req)
		if err != nil {
			e = Error(err, http.StatusBadRequest)
			goto Encode
		}
	}

	// Decode the path params.
	if h.decodePath != nil && len(p) != 0 {
		err := h.decodePath.Decode(p, req)
		if err != nil {
			e = Error(err, http.StatusBadRequest)
			goto Encode
		}
	}

	// Decode the body.
	if r.ContentLength > 0 {
		dec, err := getDecoder(h.getEncoding(r.Header, "Content-Type"))
		if err != nil {
			res = err
			goto Encode
		}

		if err := dec(r, req); err != nil {
			e = Error(err, http.StatusBadRequest)
			goto Encode
		}
	}

	res, err = h.handler(r.Context(), *req)
	if err != nil {
		e = Error(err, 0)
	}

Encode:
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")

	if e != nil {
		if !h.config.ShowPrivateErrors && e.StatusCode() == http.StatusInternalServerError {
			res = ErrInternalServerError
		} else {
			res = e
		}
	}

	if h, ok := res.(Headerer); ok {
		headers := w.Header()
		for k, v := range h.Header() {
			headers[k] = v
		}
	}

	if sc, ok := res.(StatusCoder); ok {
		w.WriteHeader(sc.StatusCode())
	} else if !h.config.DisableNoContent && h.isNil(res) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err = enc(w, res); err != nil {
		log.Println(err)
	}
}

func (h *handler[T, O]) getEncoding(header http.Header, key string) string {
	if header == nil {
		return h.config.DefaultEncoding
	}

	v := header[key]
	if len(v) == 0 {
		return h.config.DefaultEncoding
	}

	contentType := v[0]

	index := strings.Index(contentType, ";")
	if index > 0 {
		contentType = contentType[:index]
	}

	if contentType == "" || contentType == "*/*" {
		return h.config.DefaultEncoding
	}

	return strings.TrimSpace(contentType)
}

func (h *handler[T, O]) setConfig(r *Config) {
	h.config = r
	if !r.DisableNoContent {
		h.isNil = makeNilCheck(*new(O))
	}
}

const (
	headerTag = "header"
	pathTag   = "path"
	queryTag  = "query"
)
