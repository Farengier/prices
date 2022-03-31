package server

import (
	"fmt"
	"github.com/Farengier/gotools/routine"
	"github.com/rs/zerolog"
	"net/http"
)

type Server struct {
	Log zerolog.Logger
	Cfg interface {
		Addr() string
	}
}

func (s Server) Start() {
	http.HandleFunc("/", homePageHandler)

	srv := http.Server{Addr: s.Cfg.Addr()}
	routine.StartRoutine("web_server_closer", func() {
		routine.WaitTillShutdownRequested()
		srv.Close()
	})
	routine.StartRoutine("web_server", func() { srv.ListenAndServe() })
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
