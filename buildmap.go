package rebuilder

import (
	"fmt"
	"io"
)

// Buildmap maps from package name to package name. It tells you what depends
// on a given package, that is, if pkg changes it tells you what needs to be
// rebuilt
type buildMap map[string]map[string]struct{}

// add adds a dependency to the build map. If this dependency changes we need
// to rebuild this package
func (bm buildMap) add(pkg, dependency string) {
	// if dependency changes we need to change pkg
	// log.Printf(" %s -> %s", dependency, pkg)
	pkgs, ok := bm[dependency]
	if !ok {
		bm[dependency] = map[string]struct{}{
			pkg: struct{}{},
		}
		return
	}
	pkgs[pkg] = struct{}{}
}

// Rebuild tells you which packages need to be rebuilt when a package changes.
// It only tells you one level - you may need to call it repeatedly
func (bm buildMap) rebuild(pkg string) map[string]struct{} {
	return bm[pkg]
}

func (bm buildMap) show(w io.Writer) {
	for dependency, usedBys := range bm {
		fmt.Fprintf(w, "\n%s used by\n", dependency)
		for usedBy := range usedBys {
			fmt.Fprintf(w, "  %s\n", usedBy)
		}
	}
}
