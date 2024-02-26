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

// RouteGroup defines all router register functin inclide grouping function
type RouteGroup interface {
	Route
	Group() *Router
}

// Route defines all router register function
type Route interface {
	Get(path string, handlers ...Handler)
}

// Router help framework to register route path, handlers and middlewares
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

// Grouping the router
func (r *Router) Group(path string, handlers ...Handler) *Router {
	return &Router{
		basePath: path,
		handlers: handlers,
		srv:      r.srv,
	}
}

// Use adds handlers or middlewares to the router
func (r *Router) Use(handlers ...Handler) Route {
	return r
}

// Get adds route with GET HTTP method
func (r *Router) Get(path string, handlers ...Handler) {
	r.route(http.MethodGet, path, handlers...)
}

// Get adds route with POST HTTP method
func (r *Router) Post(path string, handlers ...Handler) {
	r.srv.Route(http.MethodPost, path, handlers...)
}

// Get adds route with PUT HTTP method
func (r *Router) Put(path string, handlers ...Handler) {
	r.srv.Route(http.MethodPut, path, handlers...)
}

// Get adds route with Patch HTTP method
func (r *Router) Patch(path string, handlers ...Handler) {
	r.srv.Route(http.MethodPatch, path, handlers...)
}

// Get adds route with DELETE HTTP method
func (r *Router) Delete(path string, handlers ...Handler) {
	r.srv.Route(http.MethodDelete, path, handlers...)
}

// Get adds route with mathing HTTP method that already set.
func (r *Router) Match(methods []string, path string, handlers ...Handler) {
	for _, method := range methods {
		r.srv.Route(method, path, handlers...)
	}
}

// Get adds route with any HTTP methods
func (r *Router) Any(path string, handlers ...Handler) {
	for _, method := range anyMethod {
		r.srv.Route(method, path, handlers...)
	}
}
