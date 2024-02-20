package gock

import (
	"net/http"
	"os"
	"sync"

	"github.com/chocokacang/gock/dotenv"
	"github.com/chocokacang/gock/log"
	"github.com/chocokacang/gock/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var _ Route = (*Server)(nil)

type Server struct {
	Router
	Config      *Config
	Logger      *log.Logger
	trees       trees
	maxParams   uint16
	maxSections uint16
	pool        sync.Pool
}

func New() *Server {

	dotenv.Load()

	srv := &Server{}
	srv.Config = &Config{
		APPENV:   os.Getenv("APP_ENV"),
		APPNAME:  os.Getenv("APP_NAME"),
		APPDEBUG: utils.GetEnvBool("APP_DEBUG", true),
		HTTPPORT: utils.GetEnv("HTTP_PORT", "8080"),
		HTTPH2C:  utils.GetEnvBool("HTTP_H2C", false),
		DBHOST:   os.Getenv("DB_HOST"),
		DBPORT:   os.Getenv("DB_PORT"),
		DBUSER:   os.Getenv("DB_USER"),
		DBPASS:   os.Getenv("DB_PASS"),
		LOGLEVEL: utils.GetEnv("LOG_LEVEL", "WARNING"),
		LOGFILE:  os.Getenv("LOG_FILE"),
	}
	srv.Logger = log.New("", log.LstdFlags, false, srv.Config.LOGFILE, log.ConvertLevelString(srv.Config.LOGLEVEL))
	srv.Router = Router{
		srv: srv,
	}
	srv.pool.New = func() any {
		params := make(Params, 0, srv.maxParams)
		return &ChocoKacang{params: &params, srv: srv}
	}

	srv.Logger.Debug(log.INFO, "Debug mode is enabled")

	return srv
}

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

func (srv *Server) Handler() http.Handler {
	if !srv.Config.HTTPH2C {
		return srv
	}
	srv.Logger.Debug(log.INFO, "H2C is enabled")
	h2s := &http2.Server{}
	return h2c.NewHandler(srv, h2s)
}

func (srv *Server) ServeHTTP(rsw http.ResponseWriter, rq *http.Request) {
	gock := srv.pool.Get().(*ChocoKacang)
	gock.writer.set(srv, rsw)
	gock.set(rq)

	method := gock.Request.Method
	path := gock.Request.URL.Path
	unescape := false
	if root := srv.trees[method]; root != nil {
		value := root.getValue(path, gock.params, gock.idleNodes, unescape)
		if value.params != nil {
			gock.Params = *value.params
		}
	}

	srv.pool.Put(gock)
}

func (srv *Server) Run() {
	srv.Logger.Info("Listening and serverin HTTP on port %s", srv.Config.HTTPPORT)

	server := &http.Server{
		Addr:     ":" + srv.Config.HTTPPORT,
		Handler:  srv.Handler(),
		ErrorLog: srv.Logger.WithWarningLevel(),
	}
	err := server.ListenAndServe()
	if err != nil {
		srv.Logger.Panic("%v", err)
	}
}
