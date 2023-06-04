package starlarkgrpc

import (
	_ "embed"
	"testing"

	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/stackb/grpc-starlark/pkg/protodescriptorset"
	"github.com/stackb/grpc-starlark/pkg/starlarknet"
	"github.com/stackb/grpc-starlark/pkg/starlarkthread"
)

//go:embed routeguide_proto_descriptor.pb
var routeguideProtoDescriptor []byte

var testModule = starlark.StringDict{
	"grpc":   NewModule(protoregistry.GlobalFiles),
	"net":    starlarknet.Module,
	"thread": starlarkthread.Module,
}

func RouteguideFiles(t *testing.T) *protoregistry.Files {
	pds, err := protodescriptorset.Parse(routeguideProtoDescriptor)
	if err != nil {
		t.Fatal(err)
	}
	files, err := protodesc.NewFiles(pds)
	if err != nil {
		t.Fatal(err)
	}
	return files
}
