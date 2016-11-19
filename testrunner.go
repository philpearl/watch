package rebuilder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/philpearl/rebuilder/base"
	"github.com/philpearl/rebuilder/wire"
)

type TestRunner struct {
	context *base.Context
}

func NewTestRunner(c *base.Context) *TestRunner {
	return &TestRunner{
		context: c,
	}
}

func (*TestRunner) Type() wire.TaskType { return wire.TaskTypeTest }

func (r *TestRunner) Run(pkg string) (string, error) {
	fmt.Printf("Run test for %s\n", pkg)
	dirname := filepath.Join(os.Getenv("GOPATH"), "src", pkg)

	cmd := exec.Command("go", "test")
	cmd.Dir = dirname
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Test %s failed. %v. %s\n", pkg, err, string(output))
	}
	return string(output), err
}

func (*TestRunner) WillRun(pkg string) bool { return true }
