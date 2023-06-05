package starlarkgrpc

import (
	"context"
	"fmt"
	"sort"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

// grpcClient implements starlark.Value for a grpc.Client.
type grpcClient struct {
	conn  *grpc.ClientConn
	sd    protoreflect.ServiceDescriptor
	attrs map[string]starlark.Value
}

// String implements part of the starlark.Value interface
func (*grpcClient) String() string { return "grpc.Client" }

// Type implements part of the starlark.Value interface
func (*grpcClient) Type() string { return "grpc.Client" }

// Freeze implements part of the starlark.Value interface
func (*grpcClient) Freeze() {} // immutable

// Truth implements part of the starlark.Value interface
func (*grpcClient) Truth() starlark.Bool { return starlark.True }

// Hash implements part of the starlark.Value interface
func (c *grpcClient) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

// AttrNames implements part of the starlark.HasAttrs interface
func (c *grpcClient) AttrNames() []string {
	names := make([]string, 0, len(c.attrs))
	for name := range c.attrs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Attr implements part of the starlark.HasAttrs interface
func (c *grpcClient) Attr(name string) (starlark.Value, error) {
	if attr, ok := c.attrs[name]; ok {
		return attr, nil
	} else {
		return nil, nil
	}
}

func newGrpcClient(files *protoregistry.Files) func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var serviceName string
		var channel *channel

		if err := starlark.UnpackArgs("grpc.Client", args, kwargs,
			"service", &serviceName,
			"channel", &channel,
		); err != nil {
			return nil, err
		}

		sd, ok := serviceDescriptorByName(files, serviceName)
		if !ok {
			return nil, fmt.Errorf("unknown service: %s (known: %v)", serviceName, serviceNames(files))
		}

		client := &grpcClient{
			attrs: make(map[string]starlark.Value),
			conn:  channel.ClientConn,
			sd:    sd,
		}

		methods := sd.Methods()
		for i := 0; i < methods.Len(); i++ {
			md := methods.Get(i)
			method := fmt.Sprintf("/%s/%s", sd.FullName(), md.Name())

			var attr starlark.Value
			if md.IsStreamingServer() && md.IsStreamingClient() {
				attr = starlark.NewBuiltin(method, newClientStreamingCall(method, md, client.conn))
			} else if md.IsStreamingServer() {
				attr = starlark.NewBuiltin(method, newClientStreamingCall(method, md, client.conn))
			} else if md.IsStreamingClient() {
				attr = starlark.NewBuiltin(method, newClientStreamingCall(method, md, client.conn))
			} else {
				attr = starlark.NewBuiltin(method, newClientUnaryCall(method, md, client.conn))
			}

			client.attrs[string(md.Name())] = attr
		}

		return client, nil
	}
}

func newClientUnaryCall(method string, methodDescriptor protoreflect.MethodDescriptor, conn *grpc.ClientConn) func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var in starlark.Value
		var meta md
		if err := starlark.UnpackArgs(string(methodDescriptor.Name()), args, kwargs,
			"request", &in,
			"metadata?", &meta,
		); err != nil {
			return nil, err
		}

		request, ok := protomodule.AsProtoMessage(in)
		if !ok {
			return nil, fmt.Errorf("failed to convert request argument to proto.Message: %v", in)
		}

		ctx := context.Background()

		var callOptions []grpc.CallOption
		if meta != nil {
			md := metadata.MD(meta)
			callOptions = append(callOptions, grpc.Header(&md))
			// panic(fmt.Sprintf("added headers: %v", md))
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		response := dynamicpb.NewMessage(methodDescriptor.Output())
		if err := conn.Invoke(ctx, method, request, response, callOptions...); err != nil {
			return nil, err
		}

		out, err := protomodule.NewMessage(response)
		if err != nil {
			return nil, err
		}

		return out, nil
	}
}

func newClientStreamingCall(method string, md protoreflect.MethodDescriptor, conn *grpc.ClientConn) func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := context.Background()

		clientStream, err := conn.NewStream(ctx, &grpc.StreamDesc{
			StreamName:    string(md.Name()),
			ServerStreams: md.IsStreamingServer(),
			ClientStreams: md.IsStreamingClient(),
		}, method)
		if err != nil {
			return nil, err
		}

		stream := newClientStream(clientStream, md)

		if md.IsStreamingServer() && !md.IsStreamingClient() {
			var request starlark.Value
			if err := starlark.UnpackArgs(string(md.Name()), args, kwargs,
				"request", &request,
			); err != nil {
				return nil, err
			}

			msg, ok := protomodule.AsProtoMessage(request)
			if !ok {
				return nil, fmt.Errorf("failed to convert request argument to ProtoMessage: %v", request)
			}

			if err := clientStream.SendMsg(msg); err != nil {
				return nil, err
			}

			if err := clientStream.CloseSend(); err != nil {
				return nil, err
			}
		}

		return stream, nil
	}
}
