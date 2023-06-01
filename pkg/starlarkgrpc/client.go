package starlarkgrpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"google.golang.org/grpc"
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
func (*grpcClient) Truth() starlark.Bool { return starlark.False }

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

func newClient(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var serviceName string
	var channel *channel

	if err := starlark.UnpackArgs("Client", args, kwargs,
		"service", &serviceName,
		"channel", &channel,
	); err != nil {
		return nil, err
	}

	sd, err := getServiceDescriptor(protoregistry.GlobalFiles, serviceName)
	if err != nil {
		return nil, err
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

		attrName := string(md.Name())
		var attr starlark.Value

		if md.IsStreamingServer() && md.IsStreamingClient() {
			log.Printf("TODO: grpc.Client: Registered %s (bidi stream):", attrName)
			attr = starlark.NewBuiltin(method, newClientStreamingCall(method, md, client.conn))
		} else if md.IsStreamingServer() {
			log.Printf("grpc.Client: Registered %s (server stream):", attrName)
			attr = starlark.NewBuiltin(method, newClientStreamingCall(method, md, client.conn))
		} else if md.IsStreamingClient() {
			log.Printf("TODO: grpc.Client: Registered %s (client stream):", attrName)
			attr = starlark.NewBuiltin(method, newClientStreamingCall(method, md, client.conn))
		} else {
			log.Printf("grpc.Client: Registered %s (unary method):", attrName)
			attr = starlark.NewBuiltin(method, newClientUnaryCall(method, md, client.conn))
		}

		if attr != nil {
			client.attrs[attrName] = attr
		}
	}

	return client, nil
}

func newClientUnaryCall(method string, md protoreflect.MethodDescriptor, conn *grpc.ClientConn) func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("%s requires a single argument (the request proto)", b.Name())
		}
		request, ok := protomodule.AsProtoMessage(args.Index(0))
		if !ok {
			return nil, fmt.Errorf("failed to convert request argument to ProtoMessage: %v", args.Index(0))
		}
		ctx := context.Background()
		response := dynamicpb.NewMessage(md.Output())
		if err := conn.Invoke(ctx, method, request, response); err != nil {
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
		var request starlark.Value

		if err := starlark.UnpackArgs("ClientStreamingCall", args, kwargs,
			"request", &request,
		); err != nil {
			return nil, err
		}

		msg, ok := protomodule.AsProtoMessage(request)
		if !ok {
			return nil, fmt.Errorf("failed to convert request argument to ProtoMessage: %v", request)
		}

		ctx := context.Background()

		stream, err := conn.NewStream(ctx, &grpc.StreamDesc{
			StreamName: string(md.Name()),
			// Handler:       handler.HandleStream,
			ServerStreams: md.IsStreamingServer(),
			ClientStreams: md.IsStreamingClient(),
		}, method)
		if err != nil {
			return nil, err
		}

		cstream := &clientStreamingCall{
			ClientStream: stream,
			name:         method,
			md:           md.Output(),
		}

		if err := cstream.SendMsg(msg); err != nil {
			return nil, err
		}

		if err := cstream.CloseSend(); err != nil {
			return nil, err
		}

		return cstream, nil
	}
}

type clientStreamingCall struct {
	grpc.ClientStream
	name string
	md   protoreflect.MessageDescriptor
}

func (cs *clientStreamingCall) String() string {
	return fmt.Sprintf("ClientStreamingCall<%s>", cs.name)
}
func (*clientStreamingCall) Type() string         { return "ClientStreamingCall" }
func (*clientStreamingCall) Freeze()              {} // immutable
func (*clientStreamingCall) Truth() starlark.Bool { return starlark.True }
func (c *clientStreamingCall) Iterate() starlark.Iterator {
	return &clientStreamIterator{
		ClientStream: c.ClientStream,
		md:           c.md,
	}
}

func (c *clientStreamingCall) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

func (c *clientStreamingCall) Name() string {
	return c.name
}

type clientStreamIterator struct {
	grpc.ClientStream
	md protoreflect.MessageDescriptor
}

func (it *clientStreamIterator) Next(p *starlark.Value) bool {
	msg := dynamicpb.NewMessage(it.md)
	msg.Reset()
	if err := it.ClientStream.RecvMsg(msg); err != nil {
		if err != io.EOF {
			log.Println("stream recvd error:", err)
		}
		return false
	}

	next, err := protomodule.NewMessage(msg)
	if err != nil {
		log.Println("stream message conversion error:", err)
		return false
	}
	*p = next
	return true
}

func (*clientStreamIterator) Done() {}
