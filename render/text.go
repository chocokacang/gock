package render

import (
	"fmt"
	"net/http"
)

type Text struct {
	Format string
	Data   []any
	writer http.ResponseWriter
}

func (r Text) Serve() (err error) {
	_, err = fmt.Fprintf(r.writer, "x %s", "x")
	return
}
