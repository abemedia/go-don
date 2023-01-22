package json

import (
	"github.com/abemedia/go-don"
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
	don.RegisterDecoder("application/json", decodeJSON)
	don.RegisterEncoder("application/json", encodeJSON)
}
