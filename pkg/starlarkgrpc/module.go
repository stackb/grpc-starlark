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
			"status": Status,
			"credentials": starlarkstruct.FromStringDict(
				Symbol("grpc.credentials"),
				starlark.StringDict{
					"Tls":      starlark.NewBuiltin("grpc.credentials.Tls", newTlsCredentials),
					"Insecure": starlark.NewBuiltin("grpc.credentials.Insecure", newInsecureCredentials),
				},
			),
			"Client":   starlark.NewBuiltin("grpc.Client", newGrpcClient(files)),
			"Channel":  starlark.NewBuiltin("grpc.Channel", newChannel),
			"Server":   starlark.NewBuiltin("grpc.Server", newServer(files)),
			"Error":    starlark.NewBuiltin("grpc.Error", newError),
			"Metadata": starlark.NewBuiltin("grpc.Metadata", newMetadata),
		},
	}
}
