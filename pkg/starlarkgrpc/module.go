package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "grpc",
	Members: starlark.StringDict{
		"Client":  starlark.NewBuiltin("grpc.Client", newClient),
		"Channel": starlark.NewBuiltin("grpc.Channel", newChannel),
		"Server":  starlark.NewBuiltin("grpc.Server", newServer),
		"Error":   starlark.NewBuiltin("grpc.Error", newError),
	},
}
