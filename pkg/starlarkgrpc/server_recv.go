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

// serverRecv implements starlark.Callable for the context.receive function.
type serverRecv struct {
	name string
	ss   grpc.ServerStream
	md   protoreflect.MessageDescriptor
}

func (s *serverRecv) String() string     { return "RecvRPC:" + s.name }
func (*serverRecv) Type() string         { return "RecvRPC" }
func (*serverRecv) Freeze()              {} // immutable
func (*serverRecv) Truth() starlark.Bool { return starlark.False }
func (c *serverRecv) Iterate() starlark.Iterator {
	return &serverRecvIterator{
		ss: c.ss,
		md: c.md,
	}
}

func (c *serverRecv) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

func (c *serverRecv) Name() string {
	return c.name
}

type serverRecvIterator struct {
	ss grpc.ServerStream
	md protoreflect.MessageDescriptor
}

func (it *serverRecvIterator) Next(p *starlark.Value) bool {
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

func (*serverRecvIterator) Done() {}
