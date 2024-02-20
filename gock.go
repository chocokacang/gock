package gock

import (
	"net/http"
)

type Handler func(gock *ChocoKacang)

type Handlers []Handler

type ChocoKacang struct {
	Request  *http.Request
	Response Response
	Params   Params
	params   *Params
	writer   writer
	srv      *Server
}

func (gock *ChocoKacang) set(rq *http.Request) {
	gock.Response = &gock.writer
	gock.Request = rq
}
