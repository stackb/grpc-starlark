package starlarkgrpc

import (
	"context"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var serverTransportStreamSymbol = Symbol("grpc.ServerTransportStream")

type serverTransportStream struct {
	*starlarkstruct.Struct
}

func newServerTransportStream(ctx context.Context, stream grpc.ServerTransportStream, md protoreflect.MethodDescriptor) *serverTransportStream {
	return &serverTransportStream{
		Struct: starlarkstruct.FromStringDict(
			serverTransportStreamSymbol,
			starlark.StringDict{
				"ctx":         newCtx(ctx),
				"descriptor":  newMethodDescriptor(md),
				"set_header":  starlark.NewBuiltin(string(serverTransportStreamSymbol)+".set_header", applyHeaderFunc(stream.SetHeader)),
				"send_header": starlark.NewBuiltin(string(serverTransportStreamSymbol)+".send_header", applyHeaderFunc(stream.SendHeader)),
				"set_trailer": starlark.NewBuiltin(string(serverTransportStreamSymbol)+".set_trailer", applyHeaderFunc(stream.SetTrailer)),
			},
		),
	}
}
