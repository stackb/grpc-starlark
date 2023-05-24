package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func makeStreamContext(method protoreflect.MethodDescriptor, ss grpc.ServerStream) starlark.Value {
	return starlarkstruct.FromStringDict(
		Symbol("GrpcStreamContext"),
		starlark.StringDict{
			"name":     starlark.String(method.Name()),
			"fullname": starlark.String(method.FullName()),
			"send": &sendRPC{
				name: "send",
				ss:   ss,
			},
			"recv": &recvRPC{
				name: "recv",
				ss:   ss,
				md:   method.Input(),
			},
		},
	)
}
