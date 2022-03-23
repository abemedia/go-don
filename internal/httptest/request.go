package httptest

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

func NewRequest(method, url, body string, header map[string]string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.SetRequestURI(url)

	for k, v := range header {
		ctx.Request.Header.Set(k, v)
	}

	b := []byte(body)
	ctx.Request.SetBodyStream(bytes.NewReader(b), len(b))

	return ctx
}
