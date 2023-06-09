package starlarkgrpc

import (
	"context"
	"fmt"
	"log"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// HandlerMap is a map of Handler implementations keyed by method fullname.
type HandlerMap map[string]*Handler

type HandlerRegistrationFunction func(handler *Handler) error

// Handler represents a rule implemented in starlark that implements the GrpcHandler.
type Handler struct {
	name     string
	reporter func(thread *starlark.Thread, msg string)
	// errorReporter func(msg string, args ...interface{}) error
	fn starlark.Callable
	md protoreflect.MethodDescriptor
}

func (h *Handler) Name() string {
	return h.name
}

// HandleStream implements grpc.StreamHandler for handling of server-streaming
// calls.
func (h *Handler) HandleStream(srv interface{}, ss grpc.ServerStream) error {
	var request protoreflect.ProtoMessage
	if h.md.IsStreamingServer() && !h.md.IsStreamingClient() {
		request = dynamicpb.NewMessage(h.md.Input())
		if err := ss.RecvMsg(request); err != nil {
			return err
		}
	}

	response, err := h.handle(h.md, request, ss.Context(), ss, nil)
	if err != nil {
		log.Printf("handler return value error: %v", err)
		return err
	}

	if h.md.IsStreamingClient() && !h.md.IsStreamingServer() {
		if err := ss.SendMsg(response); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) handle(method protoreflect.MethodDescriptor, request protoreflect.ProtoMessage, ctx context.Context, ss grpc.ServerStream, sts grpc.ServerTransportStream) (proto.Message, error) {
	var args starlark.Tuple

	if method.IsStreamingServer() && method.IsStreamingClient() {
		args = starlark.Tuple{newServerStream(ss, method)}
	} else if method.IsStreamingServer() {
		msg, err := protomodule.NewMessage(request)
		if err != nil {
			return nil, err
		}
		args = starlark.Tuple{newServerStream(ss, method), msg}
	} else if method.IsStreamingClient() {
		args = starlark.Tuple{newServerStream(ss, method)}
	} else {
		msg, err := protomodule.NewMessage(request)
		if err != nil {
			return nil, err
		}
		descriptor := newMethodDescriptor(method)
		args = starlark.Tuple{newServerTransportStream(ctx, sts, descriptor), msg}
	}

	thread := new(starlark.Thread)
	thread.Name = string(method.FullName())
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
