package gock

import (
	"net/http"

	"github.com/chocokacang/gock/log"
	"github.com/chocokacang/gock/render"
)

const unWritten = -1

var _ Writer = (*writer)(nil)

type Writer interface {
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

func (w *writer) RenderHeader() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *writer) Write(data []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

func (w *writer) Render(code int, r render.Render) render.Render {

	w.RenderHeader()
	return r
}
