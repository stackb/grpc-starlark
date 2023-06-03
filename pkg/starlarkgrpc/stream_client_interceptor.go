package starlarkgrpc

import (
	"context"
	"fmt"

	"go.starlark.net/starlark"
	"google.golang.org/grpc"
)

func streamClientInterceptor(parent *starlark.Thread, call starlark.Callable) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// TODO: prepare args
		args := starlark.Tuple{}

		_, err := starlark.Call(newThread(parent, call.String()), call, args, []starlark.Tuple{})
		if err != nil {
			return nil, fmt.Errorf("streamClientInterceptor error %s: %w", call.String(), err)
		}

		stream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
		return stream, nil
	}
}
