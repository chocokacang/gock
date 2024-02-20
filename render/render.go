package render

type Render interface {
	Serve() error
}

var (
	_ Render = Text{}
)
