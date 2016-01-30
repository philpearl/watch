package main

import (
	"fmt"

	"github.com/philpearl/rebuilder"
	"github.com/philpearl/rebuilder/web"
)

func main() {

	tr := rebuilder.NewTestRunner()
	w := web.NewWeb(tr)

	w.Setup()
	go w.Run()

	err := rebuilder.Watch(tr)

	fmt.Println(err.Error)
}
