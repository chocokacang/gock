package render

import (
	"fmt"
	"net/http"
)

type Text struct {
	Format string
	Data   []any
}

func (r Text) Render(w http.ResponseWriter) (err error) {
	writeContentType(w, textplain)
	_, err = fmt.Fprintf(w, r.Format, r.Data...)
	return
}
