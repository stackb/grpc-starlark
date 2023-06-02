package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func NewModule(files *protoregistry.Files) *starlarkstruct.Module {
	return &starlarkstruct.Module{
		Name: "grpc",
		Members: starlark.StringDict{
			"status":  Status,
			"Client":  starlark.NewBuiltin("grpc.Client", newClient(files)),
			"Channel": starlark.NewBuiltin("grpc.Channel", newChannel),
			"Server":  starlark.NewBuiltin("grpc.Server", newServer(files)),
			"Error":   starlark.NewBuiltin("grpc.Error", newError),
		},
	}
}
