package rebuilder

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/philpearl/rebuilder/base"
	"github.com/philpearl/rebuilder/wire"
)

type Builder struct {
	sync.Once
	context  *base.Context
	buildEnv []string
}

func NewBuilder(context *base.Context) *Builder {
	return &Builder{
		context: context,
	}
}

func (*Builder) Type() wire.TaskType { return wire.TaskTypeBuild }

func (b *Builder) Run(pkg string) (string, error) {
	var output []byte
	var err error
	dirname := filepath.Join(os.Getenv("GOPATH"), "src", pkg)

	if b.shouldBuild(dirname) {
		fmt.Printf("Build %s\n", pkg)

		cmd := exec.Command("go", "build", "-o", b.outputPath(pkg))
		cmd.Env = b.getEnv()
		cmd.Dir = dirname
		output, err = cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("Build failed. %v. %s\n", err, string(output))
		}
	}

	return string(output), err
}

func (b *Builder) WillRun(pkg string) bool {
	dirname := filepath.Join(os.Getenv("GOPATH"), "src", pkg)
	return b.shouldBuild(dirname)
}

func (b *Builder) outputPath(pkg string) string {
	// TODO: we might want to add versions here, or have a versioning builder
	return filepath.Join(b.context.BuildOutputPath, filepath.Base(pkg))
}

func (b *Builder) getEnv() []string {
	b.Once.Do(func() {
		// Environment shouldn't change while the program is running
		// env is a list of key=value
		env := os.Environ()
		eset := make(map[string]string, len(env))
		for _, val := range env {
			parts := strings.Split(val, "=")

			eset[parts[0]] = parts[1]
		}

		if b.context.BuildArch != "" {
			eset["GOARCH"] = b.context.BuildArch
		}
		if b.context.BuildOS != "" {
			eset["GOOS"] = b.context.BuildOS
		}

		env = env[0:0]
		for key, value := range eset {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}

		b.buildEnv = env
	})

	return b.buildEnv
}

// shouldBuild determines if we should build a package. We build anything that
// is a main package
func (b *Builder) shouldBuild(dirname string) bool {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirname, nil, parser.PackageClauseOnly)
	if err != nil {
		return false
	}
	for _, pkg := range pkgs {
		if pkg.Name == "main" {
			return true
		}
	}
	return false
}
