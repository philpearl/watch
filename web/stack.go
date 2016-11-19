package web

import (
	"log"
	"net/http"
	"time"

	"github.com/philpearl/rebuilder"
	"github.com/philpearl/rebuilder/base"
)

type web struct {
	track *rebuilder.Track
	mux   *http.ServeMux
}

func NewWeb(cxt *base.Context, track *rebuilder.Track) *web {
	return &web{
		track: track,
		mux:   http.NewServeMux(),
	}
}

func (w *web) Setup() {
	w.mux.HandleFunc("/status", w.status)

	w.mux.Handle("/", http.FileServer(http.Dir("public")))
}

func (w *web) Run() {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        w.mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
