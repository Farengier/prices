package renderer

import (
	"bytes"
	"github.com/Farengier/prices/src/pkg/server/deps"
	"io"
)

type template struct {
}

type page struct {
	firstHead *head
	lastHead  *head
	title     string
	content   bytes.Buffer
}

type head struct {
	data string
	next *head
}

func Template() *template {
	return &template{}
}

func (r *template) NewPage(title string) deps.Renderer {
	return &page{
		title: title,
	}
}

func (p *page) AddHead(data string) {
	h := &head{data: data}
	if p.firstHead == nil {
		p.firstHead = h
	}
	if p.lastHead == nil {
		p.lastHead = h
	} else {
		p.lastHead.next = h
	}
}

func (p *page) AddScript(script string) {
	p.AddHead("<script type=\"text/javascript\">" + script + "</script>")
}
func (p *page) AddContent(c []byte) {
	p.content.Write(c)
}

func (p *page) Render(w io.Writer) {
	w.Write([]byte("<!doctype html>"))
	w.Write([]byte("<html lang=\"en\">"))
	w.Write([]byte("<head>"))
	w.Write([]byte("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">"))
	w.Write([]byte("<link href=\"https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css\" rel=\"stylesheet\" integrity=\"sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3\" crossorigin=\"anonymous\">"))

	w.Write([]byte("<title>Prices"))
	if p.title != "" {
		w.Write([]byte(":"))
		w.Write([]byte(p.title))
	}
	w.Write([]byte("</title>"))

	for h := p.firstHead; h != nil; h = h.next {
		w.Write([]byte(h.data))
	}

	w.Write([]byte("</head>"))
	w.Write([]byte("<body>"))
	w.Write(p.content.Bytes())
	w.Write([]byte("</body>"))
}
