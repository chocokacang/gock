package gock

import (
	"net/http"
	"os"
	"sync"

	"github.com/chocokacang/gock/log"
	"github.com/chocokacang/gock/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var _ Route = (*Server)(nil)

type Server struct {
	Router
	Config *Config
	Logger *log.Logger
	pool   sync.Pool
}

func New() *Server {
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
	srv.Router = Router{}
	srv.pool.New = func() any {
		return &ChocoKacang{srv: srv}
	}
	return srv
}

func (srv *Server) Handler() http.Handler {
	if !srv.Config.HTTPH2C {
		return srv
	}
	h2s := &http2.Server{}
	return h2c.NewHandler(srv, h2s)
}

func (srv *Server) ServeHTTP(rsw http.ResponseWriter, rq *http.Request) {
	gock := srv.pool.Get().(*ChocoKacang)
	gock.writer.set(srv, rsw)
	gock.set(rq)

	gock.Response.WriteHeader(200)
	gock.Response.WriteHeader(404)
	rsw.Write([]byte("X"))
}

func (srv *Server) Run() {
	srv.Logger.Info("Run server in port :%s", srv.Config.HTTPPORT)

	server := &http.Server{
		Addr:     ":" + srv.Config.HTTPPORT,
		Handler:  srv.Handler(),
		ErrorLog: srv.Logger.WithWarningLevel(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
