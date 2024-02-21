package render

import (
	"net/http"
)

type Render interface {
	Render(writer http.ResponseWriter) error
}

var (
	_ Render = Text{}
)

var (
	textplain = []string{"text/plain; charset=utf-8"}
)

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
