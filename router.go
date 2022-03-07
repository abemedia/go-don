package don

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router interface {
	Get(path string, handle http.Handler)
	Post(path string, handle http.Handler)
	Put(path string, handle http.Handler)
	Delete(path string, handle http.Handler)
	Handle(method, path string, handle http.Handler)
	HandleFunc(method, path string, handle http.HandlerFunc)
	Group(path string) Router
	Use(mw ...Middleware)
}

type router struct {
	config      *Config
	prefix      string
	routes      []*route
	groups      []*router
	middlewares []Middleware
}

type route struct {
	method, path string
	handle       http.Handler
}

func (r *router) Get(path string, handle http.Handler) {
	r.Handle(http.MethodGet, path, handle)
}

func (r *router) Post(path string, handle http.Handler) {
	r.Handle(http.MethodPost, path, handle)
}

func (r *router) Put(path string, handle http.Handler) {
	r.Handle(http.MethodPut, path, handle)
}

func (r *router) Delete(path string, handle http.Handler) {
	r.Handle(http.MethodDelete, path, handle)
}

func (r *router) Handle(method, path string, handle http.Handler) {
	r.routes = append(r.routes, &route{method, path, withConfig(handle, r.config)})
}

func (r *router) HandleFunc(method, path string, handle http.HandlerFunc) {
	r.Handle(method, path, handle)
}

func (r *router) Group(path string) Router {
	g := &router{prefix: r.prefix + path, config: r.config}
	r.groups = append(r.groups, g)
	return g
}

func (r *router) Use(mw ...Middleware) {
	for _, fn := range mw {
		r.middlewares = append(r.middlewares, fn)
	}
}

func (r *router) addRoutes(rr *httprouter.Router, middlewares ...Middleware) {
	mw := append(middlewares, r.middlewares...)

	for _, v := range r.routes {
		h := v.handle
		for _, fn := range mw {
			h = fn(h)
		}
		rr.Handler(v.method, r.prefix+v.path, h)
	}

	for _, v := range r.groups {
		v.addRoutes(rr, mw...)
	}
}
