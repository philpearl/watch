package web

import (
	"log"
	"net/http"
	"time"

	"github.com/philpearl/rebuilder"
	"github.com/philpearl/rebuilder/base"
)

type web struct {
	testrunner *rebuilder.TestRunner
	mux        *http.ServeMux
}

func NewWeb(cxt *base.Context, testrunner *rebuilder.TestRunner) *web {
	return &web{
		testrunner: testrunner,
		mux:        http.NewServeMux(),
	}
}

func (w *web) Setup() {
	w.mux.HandleFunc("/tests", w.tests)
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
