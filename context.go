package locustr

import (
	"context"
	"time"
)

type Context struct {
	context.Context
}

func (c *Context) ReportSuccess(reqType, name string, respTime time.Duration, respLength int64) {
}

func (c *Context) ReportFailure(reqType, name string, respTime time.Duration, err error) {

}
