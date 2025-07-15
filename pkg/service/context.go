package service

import (
	"context"
	"maps"
)

type Context struct {
	reqData map[string]any
	ctx     context.Context

	cancelChan chan<- struct{}
	doneChan   chan<- struct{}
}

func (c *Context) GetRequestData() map[string]any {
	return maps.Clone(c.reqData)
}

func (c *Context) GetContext() context.Context {
	return c.ctx
}

func (c *Context) Cancel() {

}

func (c *Context) Done() {

}
