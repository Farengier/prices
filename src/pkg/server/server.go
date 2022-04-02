package server

import (
	"fmt"
	"github.com/Farengier/gotools/routine"
	"github.com/Farengier/prices/src/pkg/server/controller/admin"
	"github.com/Farengier/prices/src/pkg/server/deps"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
)

type Cfg interface {
	Addr() string
}
type Ctrl interface {
	Init()
}

type server struct {
	l  zerolog.Logger
	c  Cfg
	db interface{}
	t  deps.Template
}

func Server(l zerolog.Logger, c Cfg, db interface{}, t deps.Template) *server {
	return &server{l: l, c: c, db: db, t: t}
}

func (s *server) Start() {
	http.HandleFunc("/", homePageHandler)

	r := mux.NewRouter()
	r.HandleFunc("/", homePageHandler)
	admin.Ctrl(s.l, s.db, s.t).Init(r)

	srv := http.Server{Addr: s.c.Addr(), Handler: r}
	routine.StartRoutine("web_server_closer", func() {
		routine.WaitTillShutdownRequested()
		srv.Close()
	})
	routine.StartRoutine("web_server", func() { srv.ListenAndServe() })
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
