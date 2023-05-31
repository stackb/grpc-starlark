package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "grpc",
	Members: starlark.StringDict{
		"Server": starlark.NewBuiltin("Server", newServer),
		"Error":  starlark.NewBuiltin("Error", newError),
	},
}
