package program

import (
	"fmt"
	"log"

	"github.com/stripe/skycfg/go/protomodule"
	libproto "go.starlark.net/lib/proto"
	libtime "go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/stackb/grpc-starlark/pkg/starlarkgrpc"
	"github.com/stackb/grpc-starlark/pkg/starlarknet"
	"github.com/stackb/grpc-starlark/pkg/starlarkos"
	"github.com/stackb/grpc-starlark/pkg/starlarkthread"
	"github.com/stackb/grpc-starlark/pkg/starlarktls"
)

type Program struct {
	Files         *protoregistry.Files
	Reporter      func(msg string)
	ErrorReporter func(err error)
	Predeclared   starlark.StringDict
}

func NewProgram(files *protoregistry.Files) *Program {
	return &Program{
		Files:       files,
		Predeclared: newPredeclared(files),
		Reporter: func(msg string) {
			log.Println(msg)
		},
		ErrorReporter: func(err error) {
			fmt.Println("grpc-starlark error> ", err.Error())
		},
	}
}

func (p *Program) Init(filename string, src interface{}) (*starlark.StringDict, *starlark.Thread, error) {
	_, program, err := starlark.SourceProgram(filename, src, p.Predeclared.Has)
	if err != nil {
		return nil, nil, err
	}

	thread := new(starlark.Thread)
	thread.Name = "main"
	thread.Print = func(thread *starlark.Thread, msg string) {
		p.Reporter(msg)
	}
	libproto.SetPool(thread, p.Files)

	globals, err := program.Init(thread, p.Predeclared)
	if err != nil {
		return nil, nil, err
	}

	return &globals, thread, nil
}

func newPredeclared(files *protoregistry.Files) starlark.StringDict {
	return starlark.StringDict{
		"os":     starlarkos.Module,
		"net":    starlarknet.Module,
		"thread": starlarkthread.Module,
		"time":   libtime.Module,
		"tls":    starlarktls.Module,
		"grpc":   starlarkgrpc.NewModule(files),
		"proto":  protomodule.NewModule(fileRegistryTypes(files)),
		"struct": starlark.NewBuiltin("struct", starlarkstruct.Make),
		"module": starlark.NewBuiltin("module", starlarkstruct.MakeModule),
	}
}

func fileRegistryTypes(files *protoregistry.Files) *protoregistry.Types {
	var types protoregistry.Types
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		messages := fd.Messages()
		for i := 0; i < messages.Len(); i++ {
			md := messages.Get(i)
			msg := dynamicpb.NewMessage(md)
			msgType := msg.Type()
			types.RegisterMessage(msgType)
		}
		enums := fd.Enums()
		for i := 0; i < enums.Len(); i++ {
			ed := enums.Get(i)
			enumType := dynamicpb.NewEnumType(ed)
			types.RegisterEnum(enumType)
		}
		return true
	})
	return &types
}
