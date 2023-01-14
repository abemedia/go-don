package don

import (
	"net/http"
	"strings"

	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

type group struct {
	r      *API
	prefix string
}

func (g *group) Get(path string, handle httprouter.Handle) {
	g.Handle(http.MethodGet, path, handle)
}

func (g *group) Post(path string, handle httprouter.Handle) {
	g.Handle(http.MethodPost, path, handle)
}

func (g *group) Put(path string, handle httprouter.Handle) {
	g.Handle(http.MethodPut, path, handle)
}

func (g *group) Patch(path string, handle httprouter.Handle) {
	g.Handle(http.MethodPatch, path, handle)
}

func (g *group) Delete(path string, handle httprouter.Handle) {
	g.Handle(http.MethodDelete, path, handle)
}

func (g *group) Handle(method, path string, handle httprouter.Handle) {
	g.r.Handle(method, g.prefix+path, handle)
}

func (g *group) Handler(method, path string, handle http.Handler) {
	g.r.Handler(method, g.prefix+path, handle)
}

func (g *group) HandleFunc(method, path string, handle http.HandlerFunc) {
	g.Handler(method, path, handle)
}

func (g *group) Group(path string) Router {
	return &group{prefix: g.prefix + path, r: g.r}
}

func (g *group) Use(mw ...Middleware) {
	g.r.Use(func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		mwNext := next
		for _, fn := range mw {
			mwNext = fn(mwNext)
		}

		return func(ctx *fasthttp.RequestCtx) {
			// Only use the middleware if path belongs to group.
			if strings.HasPrefix(byteconv.Btoa(ctx.Path()), g.prefix) {
				mwNext(ctx)
			} else {
				next(ctx)
			}
		}
	})
}
