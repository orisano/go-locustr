package locustr

import (
	"context"
	"time"
)

type Context struct {
	ctx context.Context
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *Context) ReportSuccess(reqType, name string, respTime time.Duration, respLength int64) {
}

func (c *Context) ReportFailure(reqType, name string, respTime time.Duration, err error) {

}
