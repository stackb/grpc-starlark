package program

import (
	"fmt"
	"log"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/lib/proto"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"

	pkgnet "github.com/stackb/grpc-starlark/pkg/net"
	pkgos "github.com/stackb/grpc-starlark/pkg/os"
	"github.com/stackb/grpc-starlark/pkg/starlarkgrpc"
	"github.com/stackb/grpc-starlark/pkg/thread"
)

func Load(filename string, src interface{}, reporter func(msg string), errorReporter func(err error), files *protoregistry.Files) error {
	// newErrorf := func(msg string, args ...interface{}) error {
	// 	err := fmt.Errorf(filename+": "+msg, args...)
	// 	errorReporter(err)
	// 	reporter(err.Error())
	// 	return err
	// }

	predeclared := NewPredeclared(files)

	_, _, err := newProgram(filename, src, predeclared, reporter, errorReporter, files)
	if err != nil {
		return err
	}

	return nil
}

func newProgram(
	filename string,
	src interface{},
	predeclared starlark.StringDict,
	reporter func(msg string),
	errorReporter func(err error),
	files *protoregistry.Files,
) (*starlark.StringDict, *starlark.Thread, error) {
	newErrorf := func(msg string, args ...interface{}) error {
		err := fmt.Errorf(filename+": "+msg, args...)
		errorReporter(err)
		return err
	}

	_, program, err := starlark.SourceProgram(filename, src, predeclared.Has)
	if err != nil {
		return nil, nil, newErrorf("source program error: %v", err)
	}

	thread := new(starlark.Thread)
	thread.Print = func(thread *starlark.Thread, msg string) {
		reporter(msg)
	}
	proto.SetPool(thread, files)

	globals, err := program.Init(thread, predeclared)
	if err != nil {
		return nil, nil, newErrorf("eval: %w", err)
	}

	return &globals, thread, nil
}

func NewPredeclared(files *protoregistry.Files) starlark.StringDict {
	var types protoregistry.Types
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		messages := fd.Messages()
		for i := 0; i < messages.Len(); i++ {
			md := messages.Get(i)
			msg := dynamicpb.NewMessage(md)
			msgType := msg.Type()
			types.RegisterMessage(msgType)
			log.Println("Registered proto message type:", md.FullName())
		}
		enums := fd.Enums()
		for i := 0; i < enums.Len(); i++ {
			ed := enums.Get(i)
			enumType := dynamicpb.NewEnumType(ed)
			types.RegisterEnum(enumType)
			log.Println("Registered proto enum type:", ed.FullName())
		}
		return true
	})

	return starlark.StringDict{
		"os":     pkgos.Module,
		"net":    pkgnet.Module,
		"thread": thread.Module,
		"grpc":   starlarkgrpc.Module,
		"proto":  protomodule.NewModule(&types),
		"struct": starlark.NewBuiltin("struct", starlarkstruct.Make),
	}
}
