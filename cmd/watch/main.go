package main

import (
	"flag"
	"fmt"

	"github.com/philpearl/rebuilder"
	"github.com/philpearl/rebuilder/base"
	"github.com/philpearl/rebuilder/web"
)

func main() {

	c := base.NewContext()

	flag.Parse()

	tr := rebuilder.NewTestRunner(c)
	builder := rebuilder.NewBuilder(c)
	track := rebuilder.NewTrack()
	w := web.NewWeb(c, track)

	w.Setup()
	go w.Run()

	err := rebuilder.Watch(c, track, tr, builder)

	fmt.Println(err.Error)
}
