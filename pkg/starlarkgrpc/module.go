package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func NewModule(handlers HandlerMap) *starlarkstruct.Module {
	return &starlarkstruct.Module{
		Name: "grpc",
		Members: starlark.StringDict{
			"Handler": starlark.NewBuiltin("Handler", newHandlerFunction(handlers)),
			"Error":   starlark.NewBuiltin("Error", newErrorFunction()),
		},
	}
}
