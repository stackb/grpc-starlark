package starlarkgrpc

import (
	"context"
	"fmt"

	"go.starlark.net/starlark"
	"google.golang.org/grpc"
)

func unaryClientInterceptor(parent *starlark.Thread, call starlark.Callable) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// TODO: prepare args
		args := starlark.Tuple{}

		_, err := starlark.Call(newThread(parent, call.String()), call, args, []starlark.Tuple{})
		if err != nil {
			return fmt.Errorf("unaryClientInterceptor error %s: %w", call.String(), err)
		}
		return nil
	}
}

func newThread(parent *starlark.Thread, name string) *starlark.Thread {
	thread := new(starlark.Thread)
	thread.Name = name
	thread.Print = parent.Print
	thread.Load = parent.Load
	return thread
}
