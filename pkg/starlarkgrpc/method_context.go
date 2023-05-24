package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func makeMethodContext(method protoreflect.MethodDescriptor) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("GrpcMethodContext"),
		starlark.StringDict{
			"name":     starlark.String(method.Name()),
			"fullname": starlark.String(method.FullName()),
		},
	)
}
