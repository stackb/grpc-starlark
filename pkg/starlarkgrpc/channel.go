package starlarkgrpc

import (
	"fmt"

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
	var target string
	if err := starlark.UnpackArgs("Channel", args, kwargs, "target", &target); err != nil {
		return nil, err
	}
	var options []grpc.DialOption
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
