package program

import (
	"fmt"
	"log"
	"os"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/lib/proto"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/stackb/grpc-starlark/pkg/starlarkgrpc"
)

func LoadFile(filename string, reporter func(msg string), errorReporter func(err error), files *protoregistry.Files) (starlarkgrpc.HandlerMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open rule file %q: %w", filename, err)
	}
	defer f.Close()

	return Load(filename, f, reporter, errorReporter, files)
}

func Load(filename string, src interface{}, reporter func(msg string), errorReporter func(err error), files *protoregistry.Files) (starlarkgrpc.HandlerMap, error) {
	// newErrorf := func(msg string, args ...interface{}) error {
	// 	err := fmt.Errorf(filename+": "+msg, args...)
	// 	errorReporter(err)
	// 	reporter(err.Error())
	// 	return err
	// }

	handlers := make(starlarkgrpc.HandlerMap)
	predeclared := NewPredeclared(handlers, files)

	_, _, err := newProgram(filename, src, predeclared, reporter, errorReporter, files)
	if err != nil {
		return nil, err
	}

	return handlers, nil
}

func newProgram(filename string, src interface{}, predeclared starlark.StringDict, reporter func(msg string), errorReporter func(err error), files *protoregistry.Files) (*starlark.StringDict, *starlark.Thread, error) {
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

func NewPredeclared(handlers starlarkgrpc.HandlerMap, files *protoregistry.Files) starlark.StringDict {
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
		return true
	})
	return starlark.StringDict{
		"grpc":   starlarkgrpc.NewModule(handlers),
		"proto":  protomodule.NewModule(&types),
		"struct": starlark.NewBuiltin("struct", starlarkstruct.Make),
	}
}
