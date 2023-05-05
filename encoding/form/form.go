package form

import (
	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/decoder"
	"github.com/valyala/fasthttp"
)

var dec = decoder.New("form")

func decodeForm(ctx *fasthttp.RequestCtx, v any) error {
	return dec.Decode((*decoder.Args)(ctx.PostArgs()), v)
}

func decodeMultipartForm(ctx *fasthttp.RequestCtx, v any) error {
	f, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	return dec.Decode(decoder.Map(f.Value), v)
}

func init() {
	don.RegisterDecoder("application/x-www-form-urlencoded", decodeForm)
	don.RegisterDecoder("multipart/form-data", decodeMultipartForm)
}
