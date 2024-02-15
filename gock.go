package gock

import (
	"net/http"

	"github.com/chocokacang/gock/dotenv"
)

type ChocoKacang struct {
	Request  *http.Request
	Response Response
	writer   writer
	srv      *Server
}

func init() {
	dotenv.Load()
}

func (gock *ChocoKacang) set(rq *http.Request) {
	gock.Response = &gock.writer
	gock.Request = rq
}
