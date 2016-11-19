package rebuilder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/philpearl/rebuilder/base"
	"github.com/philpearl/rebuilder/wire"
	"github.com/rjeczalik/notify"
)

type Runner interface {
	Run(pkg string) (string, error)
	Type() wire.TaskType
	WillRun(pkg string) bool
}

func Watch(cxt *base.Context, track *Track, runners ...Runner) error {

	d := NewDependencies()
	fmt.Printf("Parsing dependencies\n")
	d.Load()
	fmt.Printf("Loaded\n")

	ch := make(chan notify.EventInfo, 1000)

	for _, path := range cxt.Config.WatchPaths {
		path := filepath.Join(os.Getenv("GOPATH"), "src", path, "...")
		err := notify.Watch(path, ch, notify.Write, notify.Rename, notify.Remove, notify.Create)
		if err != nil {
			return fmt.Errorf("Failed to start filesystem watch for %s. %v", path, err)
		}
	}

	paths := pathlist{}

	for ev := range ch {
		// Something has changed. Drain any queued events
		paths = paths[:0]
		paths.add(cxt, ev.Path())
	Drain:
		for {
			select {
			case ev, ok := <-ch:
				if !ok {
					break Drain
				}
				paths.add(cxt, ev.Path())
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
					if runner.WillRun(pkg) {
						track.Pending(pkg, runner.Type())
					}
				}
			}

			for _, runner := range runners {
				for pkg := range pkgs {
					track.Started(pkg, runner.Type())
					output, err := runner.Run(pkg)
					track.Ended(pkg, runner.Type(), err, output)
				}
			}
		}
	}

	return nil
}

type pathlist []string

func (l *pathlist) add(cxt *base.Context, path string) {
	// We don't want to look inside directories or files that start with a .,
	// e.g. .git
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if len(p) > 0 && p[0] == '.' {
			return
		}
	}

	// In a generic system we may want to exclude some files
	for _, excl := range cxt.Config.Skip {
		if filepath.Base(path) == excl {
			return
		}
	}

	*l = append(*l, path)
}
