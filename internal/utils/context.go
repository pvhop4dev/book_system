package utils

import (
	"context"
)

type MyContext struct {
	TraceID string
	Context context.Context
}

func NewMyContext(traceID string, ctx context.Context) *MyContext {
	return &MyContext{
		TraceID: traceID,
		Context: ctx,
	}
}
