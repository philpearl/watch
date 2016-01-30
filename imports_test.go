package rebuilder

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImports(t *testing.T) {
	imports, err := GetImportsForDir(filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "philpearl", "rebuilder"))

	sort.Strings(imports)

	assert.NoError(t, err)
	assert.Equal(t, []string{
		"fmt",
		"github.com/stretchr/testify/assert",
		"go/ast",
		"go/parser",
		"go/token",
		"io",
		"log",
		"os",
		"path/filepath",
		"sort",
		"strings",
		"testing",
	}, imports)
}
