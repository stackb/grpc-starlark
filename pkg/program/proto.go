package program

import (
	"fmt"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type skyProtoMessageType interface {
	NewMessage() protoreflect.ProtoMessage
}

func protoDecode(registry *protoregistry.Types) starlark.Callable {
	return starlark.NewBuiltin("proto.decode", func(
		t *starlark.Thread,
		fn *starlark.Builtin,
		args starlark.Tuple,
		kwargs []starlark.Tuple,
	) (starlark.Value, error) {
		var msgType starlark.Value
		var value starlark.Bytes
		if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 2, &msgType, &value); err != nil {
			return nil, err
		}
		protoMsgType, ok := msgType.(skyProtoMessageType)
		if !ok {
			return nil, fmt.Errorf("%s: for parameter 1: got %s, want proto.MessageType", fn.Name(), msgType.Type())
		}

		unmarshal := proto.UnmarshalOptions{
			Resolver: registry,
		}
		decoded := protoMsgType.NewMessage()
		if err := unmarshal.Unmarshal([]byte(value), decoded); err != nil {
			return nil, err
		}
		return protomodule.NewMessage(decoded)
	})
}

func protoEncode(registry *protoregistry.Types) starlark.Callable {
	return starlark.NewBuiltin("proto.encode", func(
		t *starlark.Thread,
		fn *starlark.Builtin,
		args starlark.Tuple,
		kwargs []starlark.Tuple,
	) (starlark.Value, error) {
		var val starlark.Value
		if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &val); err != nil {
			return nil, err
		}
		msg, ok := protomodule.AsProtoMessage(val)
		if !ok {
			return nil, fmt.Errorf("%s: for parameter 1: got %s, want proto.Message", fn.Name(), val.Type())
		}

		marshal := proto.MarshalOptions{}

		if len(kwargs) > 0 {
			if err := starlark.UnpackArgs(fn.Name(), nil, kwargs); err != nil {
				return nil, err
			}
		}

		data, err := marshal.Marshal(msg)
		if err != nil {
			return nil, err
		}

		return starlark.Bytes(data), nil
	})
}
