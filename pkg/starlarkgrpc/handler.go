package starlarkgrpc

import (
	"fmt"
	"log"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MethodHandlerMap is a map of Handler implementations keyed by method fullname.
type MethodHandlerMap map[string]*MethodHandler

type MethodHandlerRegistrationFunction func(handler *MethodHandler) error

// MethodHandler represents a rule implemented in starlark that implements the GrpcHandler.
type MethodHandler struct {
	name     string
	reporter func(thread *starlark.Thread, msg string)
	// errorReporter func(msg string, args ...interface{}) error
	fn     starlark.Callable
	method protoreflect.MethodDescriptor
}

func (h *MethodHandler) Name() string {
	return h.name
}

func (h *MethodHandler) Handle(method protoreflect.MethodDescriptor, request protoreflect.ProtoMessage, ss grpc.ServerStream) (proto.Message, error) {
	var context starlark.Value
	var args starlark.Tuple

	if method.IsStreamingServer() && method.IsStreamingClient() {
		context = makeStreamContext(method, ss)
		args = starlark.Tuple{starlark.None, context}
	} else if method.IsStreamingServer() {
		context = makeStreamContext(method, ss)
		msg, err := protomodule.NewMessage(request)
		if err != nil {
			return nil, err
		}
		args = starlark.Tuple{msg, context}
	} else if method.IsStreamingClient() {
		context = makeStreamContext(method, ss)
		args = starlark.Tuple{starlark.None, context}
	} else {
		context = makeMethodContext(method)
		msg, err := protomodule.NewMessage(request)
		if err != nil {
			return nil, err
		}
		args = starlark.Tuple{msg, context}
	}

	thread := new(starlark.Thread)
	thread.Print = h.reporter
	resp, err := starlark.Call(thread, h.fn, args, []starlark.Tuple{})
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", h.fn.String(), err)
	}

	out, ok := protomodule.AsProtoMessage(resp)
	if ok {
		return out, nil
	}

	switch t := resp.(type) {
	case *starlarkstruct.Struct:
		return nil, makeGrpcError(t)
	case starlark.NoneType:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected handler return type constructor: %v (%T)", resp, resp)
	}
}

func newRegisterHandlersFunction(callback MethodHandlerRegistrationFunction) goStarlarkFunction {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var mappings *starlark.Dict
		if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &mappings); err != nil {
			return nil, err
		}

		for _, key := range mappings.Keys() {
			name, ok := key.(starlark.String)
			if !ok {
				return nil, fmt.Errorf("%s: register error: dict key should be a fully-qualified method name (got %T)", fn.Name(), key)
			}
			value, ok, err := mappings.Get(key)
			if err != nil {
				log.Printf("registration mapping error: get %s failed: %v", key, err)
				continue
			}
			if !ok {
				panic(fmt.Sprintf("registration mapping lookup: lookup %s failed", key))
			}
			callable, ok := value.(starlark.Callable)
			if !ok {
				return nil, fmt.Errorf("%s: register error: dict value should be function (got %s)", fn.Name(), value.Type())
			}
			handler := &MethodHandler{
				name:     name.GoString(),
				fn:       callable,
				reporter: thread.Print,
			}
			callback(handler)
		}

		return starlark.None, nil
	}
}

type goStarlarkFunction func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
