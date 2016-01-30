package rebuilder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTree(t *testing.T) {

	d := NewDependencies()

	dirs := []string{
		"github.com/philpearl/tt_goji_middleware",
		"github.com/philpearl/tt_goji_middleware/base",
		"github.com/philpearl/tt_goji_middleware/postgres",
		"github.com/philpearl/tt_goji_middleware/raven",
		"github.com/philpearl/tt_goji_middleware/redis",
	}

	gopath := filepath.Join(os.Getenv("GOPATH"), "src")
	for _, dir := range dirs {
		assert.NoError(t, d.AddDir(filepath.Join(gopath, dir)))
	}

	pkgs, err := d.FindPackagesToRebuild(filepath.Join(gopath, "github.com/philpearl/tt_goji_middleware/postgres/utils.go"))
	assert.NoError(t, err)
	assert.Equal(t, map[string]struct{}{
		"github.com/philpearl/tt_goji_middleware/postgres": struct{}{},
	}, pkgs)

	pkgs, err = d.FindPackagesToRebuild(filepath.Join(gopath, "github.com/philpearl/tt_goji_middleware/raven/errorcatcher.go"))
	assert.NoError(t, err)
	assert.Equal(t, map[string]struct{}{
		"github.com/philpearl/tt_goji_middleware":       struct{}{},
		"github.com/philpearl/tt_goji_middleware/raven": struct{}{},
	}, pkgs)

}
