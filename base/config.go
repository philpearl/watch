package base

import (
	"flag"
	"fmt"
)

type Config struct {
	WatchPaths      []string
	Skip            []string
	BuildOutputPath string
	BuildArch       string
	BuildOS         string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) InitFlags() {

	flag.Var((*RepeatedStringParam)(&c.WatchPaths), "watch", "Package paths to watch for changes. Can be repeated.")
	flag.Var((*RepeatedStringParam)(&c.Skip), "skip", "Files to ignore. Can be repeated.")
	flag.StringVar(&c.BuildOutputPath, "out", "", "Path for build output.")
	flag.StringVar(&c.BuildArch, "arch", "amd64", "GOARCH value for builds.")
	flag.StringVar(&c.BuildOS, "os", "linux", "GOOS value for builds.")
}

type RepeatedStringParam []string

func (r *RepeatedStringParam) String() string {
	return fmt.Sprintf("%#v", *r)
}

func (r *RepeatedStringParam) Set(value string) error {
	*r = append(*r, value)
	return nil
}

func (r *RepeatedStringParam) Get() interface{} {
	return []string(*r)
}
