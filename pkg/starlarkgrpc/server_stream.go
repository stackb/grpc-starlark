package starlarkgrpc

import (
	"fmt"
	"io"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

var serverStreamSymbol = Symbol("grpc.ServerStream")

type serverStream struct {
	*starlarkstruct.Struct
	stream grpc.ServerStream
	md     protoreflect.MethodDescriptor
}

func (cs *serverStream) Iterate() starlark.Iterator {
	return &streamIterator{cs.stream.RecvMsg, cs.md.Input()}
}

func newServerStream(stream grpc.ServerStream, descriptor protoreflect.MethodDescriptor) *serverStream {
	return &serverStream{
		stream: stream,
		md:     descriptor,
		Struct: starlarkstruct.FromStringDict(
			serverStreamSymbol,
			starlark.StringDict{
				"ctx":         newCtx(stream.Context()),
				"descriptor":  newMethodDescriptor(descriptor),
				"set_header":  starlark.NewBuiltin(string(serverStreamSymbol)+".set_header", applyHeaderFunc(stream.SetHeader)),
				"send_header": starlark.NewBuiltin(string(serverStreamSymbol)+".send_header", applyHeaderFunc(stream.SendHeader)),
				"set_trailer": starlark.NewBuiltin(string(serverStreamSymbol)+".set_trailer", setTrailerFunc(stream.SetTrailer)),
				"send": starlark.NewBuiltin(string(serverStreamSymbol)+".send", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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
				"recv": starlark.NewBuiltin(string(serverStreamSymbol)+".recv", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					msg := dynamicpb.NewMessage(descriptor.Input())
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
			},
		),
	}
}

func applyHeaderFunc(apply func(md metadata.MD) error) func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var meta md
		if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
			"metadata", &meta,
		); err != nil {
			return nil, err
		}
		if err := apply(metadata.MD(meta)); err != nil {
			return nil, err
		}
		return starlark.None, nil
	}
}

func setTrailerFunc(apply func(md metadata.MD)) func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var md md
		if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
			"metadata", &md,
		); err != nil {
			return nil, err
		}
		apply(metadata.MD(md))
		return starlark.None, nil
	}
}
