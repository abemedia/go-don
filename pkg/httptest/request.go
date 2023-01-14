package httptest

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func NewRequest(method, url, body string, header map[string]string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.SetRequestURI(url)

	for k, v := range header {
		ctx.Request.Header.Set(k, v)
	}

	ctx.Request.SetBodyStream(strings.NewReader(body), len(body))

	return ctx
}
