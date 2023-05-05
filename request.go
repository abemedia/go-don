package don

import (
	"reflect"

	"github.com/abemedia/go-don/decoder"
	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

type requestDecoder[V any] func(v *V, ctx *fasthttp.RequestCtx, p httprouter.Params) error

func newRequestDecoder[V any](v V) requestDecoder[V] {
	query, _ := decoder.NewCached(v, "query")
	params, _ := decoder.NewCached(v, "path")
	header, _ := decoder.NewCached(v, "header")

	if query == nil && params == nil && header == nil {
		return decodeBody[V]()
	}

	return decodeRequest(query, header, params)
}

func decodeRequest[V any](query, header, params *decoder.CachedDecoder[V]) requestDecoder[V] {
	dec := decodeBody[V]()

	return func(v *V, ctx *fasthttp.RequestCtx, p httprouter.Params) error {
		val := reflect.ValueOf(v).Elem()

		if query != nil {
			if q := ctx.Request.URI().QueryArgs(); q.Len() > 0 {
				if err := query.DecodeValue((*decoder.Args)(q), val); err != nil {
					return err
				}
			}
		}

		if params != nil && len(p) > 0 {
			if err := params.DecodeValue((decoder.Params)(p), val); err != nil {
				return ErrNotFound
			}
		}

		if header != nil {
			if err := header.DecodeValue((*decoder.Header)(&ctx.Request.Header), val); err != nil {
				return err
			}
		}

		return dec(v, ctx, p)
	}
}

func decodeBody[V any]() requestDecoder[V] {
	return func(v *V, ctx *fasthttp.RequestCtx, p httprouter.Params) error {
		if ctx.Request.Header.ContentLength() == 0 {
			return nil
		}

		dec, err := getDecoder(getEncoding(ctx.Request.Header.ContentType()))
		if err != nil {
			return err
		}

		return dec(ctx, v)
	}
}
