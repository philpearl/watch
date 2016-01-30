package rebuilder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type TestResult struct {
	Err     error
	Output  string
	Started time.Time
	Ended   time.Time
	Package string
}

type TestRunner struct {
	sync.RWMutex
	results map[string]*TestResult
}

func NewTestRunner() *TestRunner {
	return &TestRunner{
		results: make(map[string]*TestResult),
	}
}

func (r *TestRunner) Run(pkg string) {
	fmt.Printf("Run test for %s\n", pkg)
	dirname := filepath.Join(os.Getenv("GOPATH"), "src", pkg)

	result := &TestResult{
		Started: time.Now(),
	}

	cmd := exec.Command("go", "test")
	cmd.Dir = dirname
	output, err := cmd.CombinedOutput()

	result.Ended = time.Now()
	result.Err = err
	result.Output = string(output)
	result.Package = pkg
	r.Lock()
	r.results[pkg] = result
	r.Unlock()

	if err != nil {
		fmt.Printf("Test %s failed. %v. %s\n", pkg, err, string(output))
	}
}

func (r *TestRunner) GetResults() []*TestResult {
	r.RLock()
	defer r.RUnlock()

	ret := make([]*TestResult, 0, len(r.results))
	for _, result := range r.results {
		if result.Err != nil {
			ret = append(ret, result)
		}
	}
	return ret
}
