package gock

type Response interface {
	Serve() error
}
