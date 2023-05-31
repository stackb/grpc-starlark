package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "grpc",
	Members: starlark.StringDict{
		"Client":  starlark.NewBuiltin("Client", newClient),
		"Channel": starlark.NewBuiltin("Channel", newChannel),
		"Server":  starlark.NewBuiltin("Server", newServer),
		"Error":   starlark.NewBuiltin("Error", newError),
	},
}
