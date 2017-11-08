package locustr

type Task struct {
	Name   string
	Weight float64
	Fn     func(ctx *Context)
}
