package starlarkgrpc

import (
	"fmt"
	"sort"

	"go.starlark.net/starlark"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type clientStream struct {
	grpc.ClientStream
	md    protoreflect.MethodDescriptor
	attrs map[string]starlark.Value
}

func newClientStream(cs grpc.ClientStream, md protoreflect.MethodDescriptor) *clientStream {
	call := &clientStream{
		ClientStream: cs,
		md:           md,
		attrs:        make(map[string]starlark.Value),
	}
	call.attrs["recv"] = &clientStreamRecv{
		recvMsg: cs.RecvMsg,
		md:      md.Output(),
	}
	if md.IsStreamingClient() {
		call.attrs["send"] = &clientStreamSend{
			sendMsg: cs.SendMsg,
			md:      md.Input(),
		}
		call.attrs["close_send"] = &clientStreamCloseSend{
			closeSend: cs.CloseSend,
		}
	}
	return call
}

func (cs *clientStream) String() string { return fmt.Sprintf("<%s %s>", cs.Type(), cs.md.FullName()) }

func (cs *clientStream) Type() string { return fmt.Sprintf("%T", cs) }

func (*clientStream) Freeze() {} // immutable

func (*clientStream) Truth() starlark.Bool { return starlark.True }

func (c *clientStream) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

// AttrNames implements part of the starlark.HasAttrs interface
func (cs *clientStream) AttrNames() (names []string) {
	for name := range cs.attrs {
		names = append(names, name)
	}
	sort.Strings(names)
	return
}

// Attr implements part of the starlark.HasAttrs interface
func (cs *clientStream) Attr(name string) (starlark.Value, error) {
	if attr, ok := cs.attrs[name]; ok {
		return attr, nil
	} else {
		return nil, nil
	}
}
