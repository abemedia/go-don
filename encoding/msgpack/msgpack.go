package msgpack

import (
	"github.com/abemedia/go-don"
	"github.com/valyala/fasthttp"
	"github.com/vmihailenco/msgpack/v5"
)

func decodeMsgpack(ctx *fasthttp.RequestCtx, v interface{}) error {
	return msgpack.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encodeMsgpack(ctx *fasthttp.RequestCtx, v interface{}) error {
	return msgpack.NewEncoder(ctx).Encode(v)
}

func init() {
	don.RegisterDecoder("application/x-msgpack", decodeMsgpack)
	don.RegisterEncoder("application/x-msgpack", encodeMsgpack)
}
