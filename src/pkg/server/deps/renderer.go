package deps

import "io"

type Template interface {
	NewPage(string) Renderer
}
type Renderer interface {
	Render(writer io.Writer)
	AddContent([]byte)
	AddScript(script string)
}
