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

var clientStreamSymbol = Symbol("grpc.ClientStream")

type clientStream struct {
	*starlarkstruct.Struct
	stream grpc.ClientStream
	md     protoreflect.MethodDescriptor
}

func (cs *clientStream) Iterate() starlark.Iterator {
	return &streamIterator{cs.stream.RecvMsg, cs.md.Output()}
}

func newClientStream(stream grpc.ClientStream, md protoreflect.MethodDescriptor) *clientStream {
	return &clientStream{
		stream: stream,
		md:     md,
		Struct: starlarkstruct.FromStringDict(
			clientStreamSymbol,
			starlark.StringDict{
				"ctx":        newCtx(stream.Context()),
				"descriptor": newMethodDescriptor(md),
				"header": starlark.NewBuiltin(string(clientStreamSymbol)+".header", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					md, err := stream.Header()
					if err != nil {
						return nil, err
					}
					return newMetadata(md), nil

				}),
				"trailer": starlark.NewBuiltin(string(clientStreamSymbol)+".trailer", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					md := stream.Trailer()
					return newMetadata(md), nil
				}),
				"recv": starlark.NewBuiltin(string(clientStreamSymbol)+".recv", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					msg := dynamicpb.NewMessage(md.Output())
					msg.Reset()
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
				"send": starlark.NewBuiltin(string(clientStreamSymbol)+".send", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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
				"close_send": starlark.NewBuiltin(string(clientStreamSymbol)+".close_send", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					if err := stream.CloseSend(); err != nil {
						return nil, err
					}
					return starlark.None, nil
				}),
			},
		),
	}
}
