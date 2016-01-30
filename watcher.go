package rebuilder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

type Runner interface {
	Run(pkg string)
}

func Watch(runners ...Runner) error {

	d := NewDependencies()
	fmt.Printf("Parsing dependencies\n")
	d.Load()
	fmt.Printf("Loaded\n")

	ch := make(chan notify.EventInfo, 1000)
	path := filepath.Join(os.Getenv("GOPATH"), "src", "github.com/unravelin", "...")
	err := notify.Watch(path, ch, notify.Write, notify.Rename, notify.Remove, notify.Create)
	if err != nil {
		return fmt.Errorf("Failed to start filesystem watch for %s. %v", path, err)
	}

	paths := pathlist{}

	for ev := range ch {
		// Something has changed. Drain any queued events
		paths = paths[:0]
		paths.add(ev.Path())
	Drain:
		for {
			select {
			case ev, ok := <-ch:
				if !ok {
					break Drain
				}
				paths.add(ev.Path())
			default:
				break Drain
			}
		}

		// Update the imports for these files
		for _, path := range paths {
			d.AddDir(path)
		}

		// Find packages that depend on the changed files
		pkgs, err := d.FindPackagesToRebuild(paths...)
		if err != nil {
			fmt.Printf("Error finding packages. %v\n", err)
		} else {
			for _, runner := range runners {
				for pkg := range pkgs {
					runner.Run(pkg)
				}
			}
		}
	}

	return nil
}

type pathlist []string

func (l *pathlist) add(path string) {
	// We don't want to look inside directories or files that start with a .,
	// e.g. .git
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if len(p) > 0 && p[0] == '.' {
			return
		}
	}

	// In a generic system we may want to exclude some files
	if filepath.Base(path) == "version.go" {
		return
	}

	*l = append(*l, path)
}
