package starlarkgrpc

import (
	"fmt"
	"math"

	"go.starlark.net/starlark"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

// channel implements starlark.Value for a grpc.ClientConn.
type channel struct {
	*grpc.ClientConn
}

// String implements part of the starlark.Value interface
func (c *channel) String() string {
	return fmt.Sprintf("grpc.Channel<%s, %v>", c.ClientConn.Target(), c.ClientConn.GetState())
}

// Type implements part of the starlark.Value interface
func (*channel) Type() string { return "grpc.Channel" }

// Freeze implements part of the starlark.Value interface
func (*channel) Freeze() {} // immutable

// Truth implements part of the starlark.Value interface
func (*channel) Truth() starlark.Bool { return starlark.False }

// Hash implements part of the starlark.Value interface
func (c *channel) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

func newChannel(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		target            string
		userAgent         string
		authority         string
		unaryInterceptor  starlark.Callable
		streamInterceptor starlark.Callable
		perRpcCredentials starlark.Callable
	)
	writeBufferSize := math.MinInt
	readBufferSize := math.MinInt
	initialWindowSize := math.MinInt
	initialConnWindowSize := math.MinInt

	if err := starlark.UnpackArgs("grpc.Channel", args, kwargs,
		"target", &target,
		"user_agent?", &userAgent,
		"authority?", &authority,
		"write_buffer_size?", &writeBufferSize,
		"read_buffer_size?", &readBufferSize,
		"initial_window_size?", &initialWindowSize,
		"initial_conn_window_size?", &initialConnWindowSize,
		"unary_interceptor?", &unaryInterceptor,
		"stream_interceptor?", &streamInterceptor,
		"per_rpc_credentials?", &perRpcCredentials,
	); err != nil {
		return nil, err
	}
	var options []grpc.DialOption
	if userAgent != "" {
		options = append(options, grpc.WithUserAgent(userAgent))
	}
	if authority != "" {
		options = append(options, grpc.WithAuthority(authority))
	}
	if writeBufferSize != math.MinInt {
		options = append(options, grpc.WithWriteBufferSize(writeBufferSize))
	}
	if readBufferSize != math.MinInt {
		options = append(options, grpc.WithWriteBufferSize(readBufferSize))
	}
	if initialWindowSize != math.MinInt {
		options = append(options, grpc.WithInitialWindowSize(int32(initialWindowSize)))
	}
	if initialConnWindowSize != math.MinInt {
		options = append(options, grpc.WithInitialWindowSize(int32(initialConnWindowSize)))
	}
	if unaryInterceptor != nil {
		options = append(options, grpc.WithUnaryInterceptor(unaryClientInterceptor(thread, unaryInterceptor)))
	}
	if streamInterceptor != nil {
		options = append(options, grpc.WithStreamInterceptor(streamClientInterceptor(thread, streamInterceptor)))
	}
	if perRpcCredentials != nil {
		options = append(options, grpc.WithPerRPCCredentials(newPerRpcCredentials(thread, perRpcCredentials)))
	}

	if len(options) == 0 {
		// options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
		options = append(options, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(target, options...)
	if err != nil {
		return nil, err
	}
	value := &channel{
		ClientConn: conn,
	}
	return value, nil
}
