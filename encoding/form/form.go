package form

import (
	"github.com/abemedia/go-don/decoder"
	"github.com/abemedia/go-don/encoding"
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
	encoding.RegisterDecoder(decodeForm, "application/x-www-form-urlencoded")
	encoding.RegisterDecoder(decodeMultipartForm, "multipart/form-data")
}
