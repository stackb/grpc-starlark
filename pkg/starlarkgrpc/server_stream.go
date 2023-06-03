package starlarkgrpc

import (
	"fmt"
	"io"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type serverStream struct {
	*starlarkstruct.Struct
	stream grpc.ServerStream
	md     protoreflect.MethodDescriptor
}

func (cs *serverStream) Iterate() starlark.Iterator {
	return &streamIterator{cs.stream.RecvMsg, cs.md.Input()}
}

func newServerStream(stream grpc.ServerStream, md protoreflect.MethodDescriptor) *serverStream {
	return &serverStream{
		stream: stream,
		md:     md,
		Struct: starlarkstruct.FromStringDict(
			Symbol("grpc.ServerStream"),
			starlark.StringDict{
				"name":                starlark.String(md.Name()),
				"fullname":            starlark.String(md.FullName()),
				"is_streaming_client": starlark.Bool(md.IsStreamingClient()),
				"is_streaming_server": starlark.Bool(md.IsStreamingServer()),
				"recv": starlark.NewBuiltin("grpc.ServerStream.recv", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					msg := dynamicpb.NewMessage(md.Input())
					if err := stream.RecvMsg(msg); err != nil {
						if err != io.EOF {
							return nil, err
						}
						return starlark.None, nil
					}
					next, err := protomodule.NewMessage(msg)
					if err != nil {
						return nil, err
					}
					return next, nil
				}),
				"send": starlark.NewBuiltin("grpc.ServerStream.send", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					for _, value := range args {
						msg, ok := protomodule.AsProtoMessage(value)
						if ok {
							if err := stream.SendMsg(msg); err != nil {
								return nil, fmt.Errorf("sending message: %w", err)
							}
						} else {
							return nil, fmt.Errorf("failed to convert send argument to ProtoMessage: %v", value)
						}
					}
					return starlark.None, nil
				}),
			},
		),
	}
}
