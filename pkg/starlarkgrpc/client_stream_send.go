package starlarkgrpc

import (
	"fmt"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type clientStreamSend struct {
	sendMsg func(interface{}) error
	md      protoreflect.MessageDescriptor
}

func (css *clientStreamSend) String() string {
	return fmt.Sprintf("<%s %s>", css.Type(), css.md.FullName())
}

func (css *clientStreamSend) Type() string { return fmt.Sprintf("%T", css) }

func (*clientStreamSend) Freeze() {} // immutable

func (*clientStreamSend) Truth() starlark.Bool { return starlark.True }

func (csr *clientStreamSend) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", csr.Type())
}

func (css *clientStreamSend) Name() string {
	return "ClientStreamSend"
}

func (css *clientStreamSend) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	for _, value := range args {
		msg, ok := protomodule.AsProtoMessage(value)
		if ok {
			if err := css.sendMsg(msg); err != nil {
				return nil, fmt.Errorf("sending message: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to convert send argument to ProtoMessage: %v", value)
		}
	}
	return starlark.None, nil
}
