package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var messageDescriptorSymbol = Symbol("protoreflect.MessageDescriptor")

type messageDescriptor struct {
	*starlarkstruct.Struct
	protoreflect.MessageDescriptor
}

func newMessageDescriptor(md protoreflect.MessageDescriptor) *messageDescriptor {
	return &messageDescriptor{
		MessageDescriptor: md,
		Struct: starlarkstruct.FromStringDict(
			messageDescriptorSymbol,
			starlark.StringDict{
				"name":     starlark.String(md.Name()),
				"fullname": starlark.String(md.FullName()),
			},
		),
	}
}
