package gock

import (
	"net/http"

	"github.com/chocokacang/gock/log"
)

const unWritten = -1

var _ Response = (*writer)(nil)

type Response interface {
	http.ResponseWriter
}

type writer struct {
	http.ResponseWriter
	size   int
	status int
	srv    *Server
}

func (w *writer) set(srv *Server, rsw http.ResponseWriter) {
	w.ResponseWriter = rsw
	w.size = unWritten
	w.status = http.StatusOK
	w.srv = srv
}

func (w *writer) Written() bool {
	return w.size != unWritten
}

func (w *writer) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			w.srv.Logger.Debug(log.WARNING, "Headers were already written. Wanted to override status code %d with %d", w.status, code)
			return
		}
		w.status = code
	}
}
