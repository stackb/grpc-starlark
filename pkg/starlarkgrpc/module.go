package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func NewModule(onRegister HandlerRegistrationFunction) *starlarkstruct.Module {
	return &starlarkstruct.Module{
		Name: "grpc",
		Members: starlark.StringDict{
			"RegisterHandlers": starlark.NewBuiltin("RegisterHandlers", newRegisterHandlersFunction(onRegister)),
			"Error":            starlark.NewBuiltin("Error", newErrorFunction()),
		},
	}
}
