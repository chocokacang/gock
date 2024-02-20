package gock

import "net/http"

var anyMethod = []string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
	http.MethodTrace,
}

var _ Route = (*Router)(nil)

type RouteGroup interface {
	Route
	Group() Route
}

type Route interface {
	Get(path string, handlers ...Handler)
}

type Router struct {
	srv *Server
}

func (r *Router) Get(path string, handlers ...Handler) {
	r.srv.Route(http.MethodGet, path, handlers...)
}

func (r *Router) Post(path string, handlers ...Handler) {
	r.srv.Route(http.MethodPost, path, handlers...)
}

func (r *Router) Put(path string, handlers ...Handler) {
	r.srv.Route(http.MethodPut, path, handlers...)
}

func (r *Router) Patch(path string, handlers ...Handler) {
	r.srv.Route(http.MethodPatch, path, handlers...)
}

func (r *Router) Delete(path string, handlers ...Handler) {
	r.srv.Route(http.MethodDelete, path, handlers...)
}

func (r *Router) Match(methods []string, path string, handlers ...Handler) {
	for _, method := range methods {
		r.srv.Route(method, path, handlers...)
	}
}

func (r *Router) Any(path string, handlers ...Handler) {
	for _, method := range anyMethod {
		r.srv.Route(method, path, handlers...)
	}
}
