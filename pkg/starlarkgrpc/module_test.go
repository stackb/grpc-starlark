package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/stackb/grpc-starlark/pkg/starlarknet"
	"github.com/stackb/grpc-starlark/pkg/starlarkthread"
)

var testModule = starlark.StringDict{
	"grpc":   NewModule(protoregistry.GlobalFiles),
	"net":    starlarknet.Module,
	"thread": starlarkthread.Module,
}
