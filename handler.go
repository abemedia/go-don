package don

import (
	"bytes"
	"context"
	"net/http"

	"github.com/abemedia/go-don/decoder"
	"github.com/abemedia/go-don/internal"
	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

// StatusCoder allows you to customise the HTTP response code.
type StatusCoder interface {
	StatusCode() int
}

// Headerer allows you to customise the HTTP headers.
type Headerer interface {
	Header() http.Header
}

// Handle is the type for your handlers.
type Handle[T any, O any] func(ctx context.Context, request T) (O, error)

// H wraps your handler function with the Go generics magic.
//nolint:gocognit,cyclop
func H[T any, O any](handle Handle[T, O]) httprouter.Handle {
	var (
		decodeHeader *decoder.HeaderDecoder
		decodePath   *decoder.ParamsDecoder
		decodeQuery  *decoder.ArgsDecoder
		isNil        = makeNilCheck(*new(O))
	)

	{
		var t T
		if hasTag(t, headerTag) {
			dec, err := decoder.NewHeaderDecoder(t, headerTag)
			if err == nil {
				decodeHeader = dec
			}
		}
		if hasTag(t, queryTag) {
			dec, err := decoder.NewArgsDecoder(t, queryTag)
			if err == nil {
				decodeQuery = dec
			}
		}
		if hasTag(t, pathTag) {
			dec, err := decoder.NewParamsDecoder(t, pathTag)
			if err == nil {
				decodePath = dec
			}
		}
	}

	return func(ctx *fasthttp.RequestCtx, p httprouter.Params) {
		contentType := getEncoding(ctx.Request.Header.Peek(fasthttp.HeaderAccept))

		enc, err := getEncoder(contentType)
		if err != nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotAcceptable), fasthttp.StatusNotAcceptable)
			return
		}

		req := new(T)

		var (
			res any
			e   *HTTPError
		)

		// Decode the header.
		if decodeHeader != nil {
			err = decodeHeader.Decode(&ctx.Request.Header, req)
			if err != nil {
				e = Error(err, http.StatusBadRequest)
				goto Encode
			}
		}

		// Decode the URL query.
		if decodeQuery != nil {
			if q := ctx.URI().QueryArgs(); q.Len() > 0 {
				err := decodeQuery.Decode(q, req)
				if err != nil {
					e = Error(err, http.StatusBadRequest)
					goto Encode
				}
			}
		}

		// Decode the path params.
		if decodePath != nil && len(p) != 0 {
			err := decodePath.Decode(p, req)
			if err != nil {
				e = Error(err, http.StatusBadRequest)
				goto Encode
			}
		}

		// Decode the body.
		if ctx.Request.Header.ContentLength() > 0 {
			dec, err := getDecoder(getEncoding(ctx.Request.Header.ContentType()))
			if err != nil {
				res = err
				goto Encode
			}

			if err := dec(ctx, req); err != nil {
				e = Error(err, http.StatusBadRequest)
				goto Encode
			}
		}

		res, err = handle(ctx, *req)
		if err != nil {
			e = Error(err, 0)
		}

	Encode:
		ctx.SetContentType(contentType + "; charset=utf-8")

		if e != nil {
			res = e
		}

		if h, ok := res.(Headerer); ok {
			for k, v := range h.Header() {
				ctx.Response.Header.Set(k, v[0])
			}
		}

		if sc, ok := res.(StatusCoder); ok {
			ctx.SetStatusCode(sc.StatusCode())
		}

		if isNil(res) {
			res = nil
			ctx.Response.Header.SetContentLength(-3)
		}

		if err = enc(ctx, res); err != nil {
			ctx.Logger().Printf("%v", err)
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		}
	}
}

func getEncoding(b []byte) string {
	index := bytes.IndexRune(b, ';')
	if index > 0 {
		b = b[:index]
	}

	return internal.Btoa(bytes.TrimSpace(b))
}

const (
	headerTag = "header"
	pathTag   = "path"
	queryTag  = "query"
)
