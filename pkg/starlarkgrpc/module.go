package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func NewModule(onRegister MethodHandlerRegistrationFunction) *starlarkstruct.Module {
	return &starlarkstruct.Module{
		Name: "grpc",
		Members: starlark.StringDict{
			"RegisterHandlers": starlark.NewBuiltin("RegisterHandlers", newRegisterHandlersFunction(onRegister)),
			"Server":           starlark.NewBuiltin("Server", newServerFunction()),
			"Listener":         starlark.NewBuiltin("Listener", newNetListenerFunction()),
			"Error":            starlark.NewBuiltin("Error", newErrorFunction()),
		},
	}
}
