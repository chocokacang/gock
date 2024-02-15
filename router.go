package gock

var _ Route = (*Router)(nil)

type RouteGroup interface {
	Route
	Group() Route
}

type Route interface {
	Get()
}

type Router struct {
	srv *Server
}

func (r *Router) Get() {

}
