package base

type Context struct {
	*Config
}

func NewContext() *Context {
	cxt := &Context{}
	cxt.Config = NewConfig()
	cxt.Config.InitFlags()

	return cxt
}
