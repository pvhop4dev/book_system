package utils

import (
	"context"
	"time"
)

type MyContext struct {
	Context context.Context
	TraceID string
}

func NewMyContext(ctx context.Context, traceID string) *MyContext {
	return &MyContext{
		TraceID: traceID,
		Context: ctx,
	}
}

func (c *MyContext) Deadline() (deadline time.Time, ok bool) {
	return c.Context.Deadline()
}

func (c *MyContext) Done() <-chan struct{} {
	return c.Context.Done()
}

func (c *MyContext) Err() error {
	return c.Context.Err()
}

func (c *MyContext) Value(key any) any {
	return c.Context.Value(key)
}
