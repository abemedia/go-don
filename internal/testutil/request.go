package testutil

import (
	"strings"

	"github.com/valyala/fasthttp"
)

// NewRequest returns a new [fasthttp.RequestCtx] with the given method, url, body and header.
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
