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

// HandlerMap is a map of Handler implementations keyed by method fullname.
type HandlerMap map[string]*Handler

// Handler represents a rule implemented in starlark that implements the GrpcHandler.
type Handler struct {
	name     string
	reporter func(thread *starlark.Thread, msg string)
	// errorReporter func(msg string, args ...interface{}) error
	handler *starlarkstruct.Struct
}

func (h *Handler) Name() string {
	val, err := h.handler.Attr("name")
	if err != nil {
		log.Fatalf(".name access error: %v", err)
	}
	return val.(*starlark.String).GoString()
}

func (h *Handler) Handle(method protoreflect.MethodDescriptor, request protoreflect.ProtoMessage, ss grpc.ServerStream) (proto.Message, error) {
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

	val, err := h.handler.Attr("impl")
	if err != nil {
		return nil, fmt.Errorf(".impl access error: %w", err)
	}
	callable := val.(starlark.Callable)

	thread := new(starlark.Thread)
	thread.Print = h.reporter
	resp, err := starlark.Call(thread, callable, args, []starlark.Tuple{})
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", callable.String(), err)
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

func newHandlerFunction(handlers HandlerMap) goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var implementation starlark.Callable

		if err := starlark.UnpackArgs("Handler", args, kwargs,
			"name", &name,
			"impl", &implementation,
		); err != nil {
			return nil, err
		}

		handler := starlarkstruct.FromStringDict(
			Symbol("Handler"),
			starlark.StringDict{
				"name": starlark.String(name),
				"impl": implementation,
			},
		)

		handlers[name] = &Handler{
			name:     name,
			handler:  handler,
			reporter: thread.Print,
			// errorReporter: newErrorf,
		}

		return handler, nil
	}
}

type goStarlarkFunction func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
