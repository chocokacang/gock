package gock

import (
	"net/http"

	"github.com/chocokacang/gock/db"
	"github.com/chocokacang/gock/render"
)

type Handler func(gock *ChocoKacang) Response

type Handlers []Handler

type ChocoKacang struct {
	Request   *http.Request
	Writer    Writer
	Params    Params
	params    *Params
	writer    writer
	handlers  Handlers
	index     int8
	srv       *Server
	idleNodes *[]idleNode
	fullPath  string
}

func (gock *ChocoKacang) set(rq *http.Request) {
	gock.Writer = &gock.writer
	gock.Request = rq
	gock.index = -1

	*gock.params = (*gock.params)[:0]
	*gock.idleNodes = (*gock.idleNodes)[:0]
}

func (gock *ChocoKacang) Next() {
	gock.index++
	for gock.index < int8(len(gock.handlers)) {
		if gock.Writer.Written() {
			return
		}
		gock.handlers[gock.index](gock).Next()
	}
	gock.writer.RenderHeader()
}

func (gock *ChocoKacang) Render(code int, r render.Render) Response {
	gock.Writer.WriteHeader(code)

	err := r.Render(gock.Writer)
	if err != nil {
		gock.srv.Logger.Error("%v", err)
	}

	gock.writer.RenderHeader()
	return gock
}

func (gock *ChocoKacang) Text(code int, format string, v ...any) Response {
	return gock.Render(code, render.Text{Format: format, Data: v})
}

func (gock *ChocoKacang) SQL(query string) db.Query {
	return gock.srv.db.Query(query)
}

func (gock *ChocoKacang) Model(dst interface{}) db.ModelStatement {
	model := &db.Model{Dst: dst}
	return model.DB(gock.srv.db)
}
