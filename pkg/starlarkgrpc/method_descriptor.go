package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var methodDescriptorSymbol = Symbol("protoreflect.MethodDescriptor")

type methodDescriptor struct {
	*starlarkstruct.Struct
	protoreflect.MethodDescriptor
}

func newMethodDescriptor(md protoreflect.MethodDescriptor) *methodDescriptor {
	return &methodDescriptor{
		MethodDescriptor: md,
		Struct: starlarkstruct.FromStringDict(
			methodDescriptorSymbol,
			starlark.StringDict{
				"name":                starlark.String(md.Name()),
				"fullname":            starlark.String(md.FullName()),
				"is_streaming_client": starlark.Bool(md.IsStreamingClient()),
				"is_streaming_server": starlark.Bool(md.IsStreamingServer()),
				"input":               newMessageDescriptor(md.Input()),
				"output":              newMessageDescriptor(md.Output()),
			},
		),
	}
}
