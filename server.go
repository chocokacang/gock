package gock

import (
	"net/http"
	"sync"

	"github.com/chocokacang/gock/config"
	"github.com/chocokacang/gock/db"
	"github.com/chocokacang/gock/dotenv"
	"github.com/chocokacang/gock/log"
	"github.com/chocokacang/gock/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var _ Route = (*Server)(nil)

// Server is a framework instance contains configuration, router and logger instance.
// Create an instance of framework, by New()
type Server struct {
	Router
	Config      config.Config
	Logger      *log.Logger
	trees       trees
	maxParams   uint16
	maxSections uint16
	pool        sync.Pool
	db          *db.DB
}

// New returns a framework instance without any middleware attached
func New() *Server {
	dotenv.Load()

	config := config.Default()
	logger := log.Default(config.LogLevel, config.LogFile)

	srv := &Server{
		Config: config,
		Logger: logger,
	}

	srv.Router = Router{
		basePath: "/",
		srv:      srv,
	}

	srv.db = &db.DB{
		Logger: srv.Logger,
	}

	srv.pool.New = func() any {
		params := make(Params, 0, srv.maxParams)
		idleNodes := make([]idleNode, 0, srv.maxSections)
		return &ChocoKacang{params: &params, srv: srv, idleNodes: &idleNodes}
	}

	srv.Logger.Debug(log.INFO, "Debug mode is enabled. The log level automatically set to INFO level.")

	return srv
}

func (srv *Server) OpenDB(name, driver, dsn string) {
	srv.db.Open(name, driver, dsn)
}

// Route register new route path to the framework
func (srv *Server) Route(method, path string, handlers ...Handler) {
	if method == "" {
		srv.Logger.Panic("HTTP method can not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		srv.Logger.Panic("Route path must begin with \"/\"")
	}
	if len(handlers) < 1 {
		srv.Logger.Panic("Route must have at lease one")
	}

	srv.Logger.Debug(log.INFO, "Add Route: %s \"%s\" (%s handlers)", method, path, utils.GetFunctionName(handlers[len(handlers)-1]))

	if srv.trees == nil {
		srv.trees = make(trees)
	}

	root := srv.trees[method]
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		srv.trees[method] = root
	}

	root.addRoute(path, handlers)

	if paramsCount := countParams(path); paramsCount > srv.maxParams {
		srv.maxParams = paramsCount
	}

	if sectionsCount := countSections(path); sectionsCount > srv.maxSections {
		srv.maxSections = sectionsCount
	}
}

// Handler returns http.Handler
func (srv *Server) Handler() http.Handler {
	if !srv.Config.HTTPH2C {
		return srv
	}
	srv.Logger.Debug(log.INFO, "H2C is enabled")
	h2s := &http2.Server{}
	return h2c.NewHandler(srv, h2s)
}

func (srv *Server) handle(gock *ChocoKacang) {
	method := gock.Request.Method
	path := gock.Request.URL.Path
	unescape := true
	if root := srv.trees[method]; root != nil {
		value := root.getValue(path, gock.params, gock.idleNodes, unescape)
		if value.params != nil {
			gock.Params = *value.params
		}
		if value.handlers != nil {
			gock.handlers = value.handlers
			gock.fullPath = value.fullPath
			gock.Next()
			srv.Logger.Info("%s %s %s %d", method, path, gock.Request.Proto, gock.Writer.Status())
			return
		}
	}
}

// ServeHTTP use for handle HTTP request from the clients
func (srv *Server) ServeHTTP(rsw http.ResponseWriter, rq *http.Request) {
	gock := srv.pool.Get().(*ChocoKacang)
	gock.writer.set(srv, rsw)
	gock.set(rq)

	srv.handle(gock)

	srv.pool.Put(gock)
}

// Run web server framework
func (srv *Server) Run() {
	srv.Logger.Info("Listening and serving HTTP on port %s", srv.Config.HTTPPort)

	server := &http.Server{
		Addr:     ":" + srv.Config.HTTPPort,
		Handler:  srv.Handler(),
		ErrorLog: srv.Logger.WithErrorLevel(),
	}
	err := server.ListenAndServe()
	if err != nil {
		srv.Logger.Panic("%v", err)
	}
}
