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
	w := web.NewWeb(c, tr) // TODO: add builder, or a results store

	w.Setup()
	go w.Run()

	err := rebuilder.Watch(c, tr, builder)

	fmt.Println(err.Error)
}
