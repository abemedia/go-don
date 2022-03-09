package don

import (
	"net/http"
	"strings"
)

type group struct {
	r      *API
	prefix string
}

func (g *group) Get(path string, handle http.Handler) {
	g.Handle(http.MethodGet, path, handle)
}

func (g *group) Post(path string, handle http.Handler) {
	g.Handle(http.MethodPost, path, handle)
}

func (g *group) Put(path string, handle http.Handler) {
	g.Handle(http.MethodPut, path, handle)
}

func (g *group) Patch(path string, handle http.Handler) {
	g.Handle(http.MethodPatch, path, handle)
}

func (g *group) Delete(path string, handle http.Handler) {
	g.Handle(http.MethodDelete, path, handle)
}

func (g *group) Handle(method, path string, handle http.Handler) {
	g.r.Handle(method, g.prefix+path, handle)
}

func (g *group) HandleFunc(method, path string, handle http.HandlerFunc) {
	g.Handle(method, path, handle)
}

func (g *group) Group(path string) Router {
	return &group{prefix: g.prefix + path, r: g.r}
}

func (g *group) Use(mw ...Middleware) {
	g.r.Use(func(next http.Handler) http.Handler {
		mwNext := next
		for _, fn := range mw {
			mwNext = fn(mwNext)
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only use the middleware if path belongs to group.
			if strings.HasPrefix(r.URL.Path, g.prefix) {
				mwNext.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	})
}
