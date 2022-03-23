package form

import (
	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/decoder"
	"github.com/valyala/fasthttp"
)

var dec = decoder.NewDecoder("form")

func decodeForm(ctx *fasthttp.RequestCtx, v interface{}) error {
	return dec.Decode((*decoder.ArgsGetter)(ctx.PostArgs()), v)
}

func decodeMultipartForm(ctx *fasthttp.RequestCtx, v interface{}) error {
	f, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	return dec.Decode(decoder.MapGetter(f.Value), v)
}

func init() {
	don.RegisterDecoder("application/x-www-form-urlencoded", decodeForm)
	don.RegisterDecoder("multipart/form-data", decodeMultipartForm)
}
