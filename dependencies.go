package rebuilder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

/*

$GOPATH/src/github.com/philpearl/ut
- package is github.com/philpearl/ut
- imports bytes, fmt, etc

$GOPATH/src/github.com/philpearl/ut/example
- package is example
imports github.com/philpearl/ut

If $GOPATH/src/github.com/philpearl/ut changes
1 find the package that directory represents
2 if not "main", find directories that include that package (and emit the directory & package name)
3 goto 1

Note the package name imported implies a search path, and the package name the
package declares does not have to match. We do not need to concern ourselves
with the internal package name
*/

type dependencies struct {
	gopath      string
	pkgImports  map[string][]string
	pkgsToBuild buildMap
}

func NewDependencies() *dependencies {
	return &dependencies{
		gopath:      filepath.Join(os.Getenv("GOPATH"), "src"),
		pkgImports:  make(map[string][]string),
		pkgsToBuild: make(buildMap),
	}
}

// AddDir adds a directory to the dependency map
func (d *dependencies) AddDir(dirname string) error {
	// Get the package name this directory represents. Note we only want the
	// name the filesystem heirachy implies - the package may declare another
	// name but that is not relevant to us
	pkg, err := d.dirnameToPkgName(dirname)
	if err != nil {
		return err
	}

	// Get what this package depends on
	imports, err := GetImportsForDir(dirname)
	if err != nil {
		return err
	}

	d.pkgImports[pkg] = imports

	// Add to the reverse map: each of these imports should point to this package.
	// The package should always point to itself in this map
	d.pkgsToBuild.add(pkg, pkg)
	for _, animport := range imports {
		d.pkgsToBuild.add(pkg, animport)
	}
	return nil
}

// FindPackagesToRebuild finds all packages that may need to be rebuilt if the
// given file has changed
func (d *dependencies) FindPackagesToRebuild(filename string) (map[string]struct{}, error) {
	pkgname, err := d.dirnameToPkgName(filepath.Dir(filename))
	if err != nil {
		return nil, err
	}

	// We know we need to rebuild the current package.
	pkgs := map[string]struct{}{
		pkgname: struct{}{},
	}
	// Add packages that depend on the current one.
	d.getBuildTree(pkgs, pkgname)

	return pkgs, nil
}

// getBuildTree adds packages to build to pkgs which are dependent on currPkg
func (d *dependencies) getBuildTree(pkgs map[string]struct{}, currPkg string) {
	// Find packages which are immediately dependent on currPkg
	newPkgs := d.pkgsToBuild.rebuild(currPkg)
	for pkg := range newPkgs {
		// If we've not seen this package yet then add it to our set of packages,
		// and recurse down to ensure we have packages that depend on it.
		if _, ok := pkgs[pkg]; !ok {
			pkgs[pkg] = struct{}{}
			d.getBuildTree(pkgs, pkg)
		}
	}
	return
}

// Load reads the gopath and builds a complete dependency map
func (d *dependencies) Load() error {
	err := filepath.Walk(d.gopath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name()[0] == '.' {
				// Directories starting with . should be ignored, and we don't
				// want to decend into them
				return filepath.SkipDir
			}
			err := d.AddDir(path)
			if err != nil {
				// We'll carry on despite any errors
				log.Println(err.Error())
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error %v\n", err)
	}

	return err
}

func (d *dependencies) dirnameToPkgName(dir string) (string, error) {
	pkgName, err := filepath.Rel(d.gopath, dir)
	if err != nil {
		return "", fmt.Errorf("Cannot determine package for directory %s. %v", dir, err)
	}

	return pkgName, err
}
