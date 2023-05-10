package msgpack

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/valyala/fasthttp"
	"github.com/vmihailenco/msgpack/v5"
)

func decodeMsgpack(ctx *fasthttp.RequestCtx, v any) error {
	return msgpack.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encodeMsgpack(ctx *fasthttp.RequestCtx, v any) error {
	return msgpack.NewEncoder(ctx).Encode(v)
}

func init() {
	encoding.RegisterDecoder(decodeMsgpack, "application/x-msgpack")
	encoding.RegisterEncoder(encodeMsgpack, "application/x-msgpack")
}
