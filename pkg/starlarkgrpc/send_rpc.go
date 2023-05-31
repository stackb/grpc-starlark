package starlarkgrpc

import (
	"fmt"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"google.golang.org/grpc"
)

// sendRPC implements starlark.Callable for the context.send function.
type sendRPC struct {
	name string
	ss   grpc.ServerStream
}

func (r *sendRPC) String() string     { return r.name }
func (*sendRPC) Type() string         { return "SendRPC" }
func (*sendRPC) Freeze()              {} // immutable
func (*sendRPC) Truth() starlark.Bool { return starlark.False }
func (c *sendRPC) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

func (c *sendRPC) Name() string {
	return c.name
}

func (c *sendRPC) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	for _, value := range args {
		msg, ok := protomodule.AsProtoMessage(value)
		if ok {
			if err := c.ss.SendMsg(msg); err != nil {
				return nil, fmt.Errorf("sending message: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to convert send argument to ProtoMessage: %v", value)
		}
	}
	return starlark.None, nil
}
