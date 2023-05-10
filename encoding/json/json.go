package json

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

func decodeJSON(ctx *fasthttp.RequestCtx, v any) error {
	return json.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encodeJSON(ctx *fasthttp.RequestCtx, v any) error {
	return json.NewEncoder(ctx).Encode(v)
}

func init() {
	encoding.RegisterDecoder(decodeJSON, "application/json")
	encoding.RegisterEncoder(encodeJSON, "application/json")
}
