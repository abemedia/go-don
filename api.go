package don

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Empty struct{}

var DefaultEncoding = "text/plain"

type Middleware func(http.Handler) http.Handler

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

type API struct {
	router *httprouter.Router
	config *Config
	mw     []Middleware
}

type Config struct {
	DefaultEncoding string
}

// New creates a new API instance.
func New(c *Config) *API {
	if c == nil {
		c = &Config{}
	}
	if c.DefaultEncoding == "" {
		c.DefaultEncoding = DefaultEncoding
	}

	r := httprouter.New()
	r.NotFound = withConfig(E(ErrNotFound), c)
	r.MethodNotAllowed = withConfig(E(ErrMethodNotAllowed), c)

	return &API{router: r, config: c}
}

// Router creates a http.Handler for the API.
func (r *API) Router() http.Handler {
	h := http.Handler(r.router)
	for _, mw := range r.mw {
		h = mw(h)
	}
	return h
}

func (r *API) Get(path string, handle http.Handler) {
	r.Handle(http.MethodGet, path, handle)
}

func (r *API) Post(path string, handle http.Handler) {
	r.Handle(http.MethodPost, path, handle)
}

func (r *API) Put(path string, handle http.Handler) {
	r.Handle(http.MethodPut, path, handle)
}

func (r *API) Patch(path string, handle http.Handler) {
	r.Handle(http.MethodPatch, path, handle)
}

func (r *API) Delete(path string, handle http.Handler) {
	r.Handle(http.MethodDelete, path, handle)
}

func (r *API) Handle(method, path string, handle http.Handler) {
	var hh httprouter.Handle
	if h, ok := handle.(Handler); ok {
		hh = withConfig(h, r.config).handle
	} else {
		hh = wrapHandler(handle)
	}
	r.router.Handle(method, path, hh)
}

func (r *API) HandleFunc(method, path string, handle http.HandlerFunc) {
	r.Handle(method, path, handle)
}

func (r *API) Group(path string) Router {
	return &group{prefix: path, r: r}
}

func (r *API) Use(mw ...Middleware) {
	r.mw = append(r.mw, mw...)
}

func withConfig(handle Handler, c *Config) Handler {
	if h, ok := handle.(interface{ setConfig(*Config) }); ok {
		h.setConfig(c)
	}
	return handle
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}
