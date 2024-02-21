package gock

import (
	"net/http"

	"github.com/chocokacang/gock/utils"
)

var anyMethod = []string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
	http.MethodTrace,
}

var _ Route = (*Router)(nil)

type RouteGroup interface {
	Route
	Group() *Router
}

type Route interface {
	Get(path string, handlers ...Handler)
}

type Router struct {
	basePath string
	handlers Handlers
	srv      *Server
}

func (r *Router) route(method string, path string, handlers ...Handler) {
	finalPath := utils.JointPath(r.basePath, path)
	finalHandlers := r.combineHandlers(handlers)
	r.srv.Route(method, finalPath, finalHandlers...)
}

func (r *Router) combineHandlers(handlers Handlers) Handlers {
	oldSize := len(r.handlers)
	newSize := len(handlers)
	finalSize := oldSize + newSize
	finalHandlers := make(Handlers, finalSize)
	copy(finalHandlers, r.handlers)
	copy(finalHandlers[oldSize:], handlers)
	return finalHandlers
}

func (r *Router) Group(path string, handlers ...Handler) *Router {
	return &Router{
		basePath: path,
		handlers: handlers,
		srv:      r.srv,
	}
}

func (r *Router) Use(handlers ...Handler) Route {
	return r
}

func (r *Router) Get(path string, handlers ...Handler) {
	r.route(http.MethodGet, path, handlers...)
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
