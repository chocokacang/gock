package gock

import (
	"net/http"

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
	idleNodes *[]idleNode
	srv       *Server
}

func (gock *ChocoKacang) set(rq *http.Request) {
	gock.Writer = &gock.writer
	gock.Request = rq
}

func (gock *ChocoKacang) Response(code int, r render.Render) Response {
	return gock.writer.Render(code, r)
}

func (gock *ChocoKacang) Text(code int, format string, v ...any) Response {
	return gock.Response(code, render.Text{Format: format, Data: v})
}
