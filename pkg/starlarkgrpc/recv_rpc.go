package starlarkgrpc

import (
	"fmt"
	"io"
	"log"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// recvRPC implements starlark.Callable for the context.receive function.
type recvRPC struct {
	name string
	ss   grpc.ServerStream
	md   protoreflect.MessageDescriptor
}

func (*recvRPC) String() string       { return "RecvRPC" }
func (*recvRPC) Type() string         { return "RecvRPC" }
func (*recvRPC) Freeze()              {} // immutable
func (*recvRPC) Truth() starlark.Bool { return starlark.False }
func (c *recvRPC) Iterate() starlark.Iterator {
	return &recvRpcIterator{
		ss: c.ss,
		md: c.md,
	}
}

func (c *recvRPC) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

func (c *recvRPC) Name() string {
	return c.name
}

type recvRpcIterator struct {
	ss grpc.ServerStream
	md protoreflect.MessageDescriptor
}

func (it *recvRpcIterator) Next(p *starlark.Value) bool {
	msg := dynamicpb.NewMessage(it.md)
	msg.Reset()
	if err := it.ss.RecvMsg(msg); err != nil {
		if err != io.EOF {
			log.Println("stream recvd error:", err)
		}
		return false
	}

	next, err := protomodule.NewMessage(msg)
	if err != nil {
		log.Println("stream message conversion error:", err)
		return false
	}
	*p = next
	return true
}

func (*recvRpcIterator) Done() {}
