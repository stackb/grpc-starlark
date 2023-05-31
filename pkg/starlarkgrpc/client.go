package starlarkgrpc

import (
	"context"
	"fmt"
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
		key := string(md.Name())
		var val starlark.Value
		if md.IsStreamingServer() && md.IsStreamingClient() {
			log.Printf("TODO: grpc.Client: Registered %s (bidi stream):", key)
		} else if md.IsStreamingServer() {
			log.Printf("TODO: grpc.Client: Registered %s (server stream):", key)
		} else if md.IsStreamingClient() {
			log.Printf("TODO: grpc.Client: Registered %s (client stream):", key)
		} else {
			log.Printf("grpc.Client: Registered %s (unary method):", key)
			name := fmt.Sprintf("/%s/%s", sd.FullName(), md.Name())
			call := &unaryCall{
				name: name,
				md:   md,
				conn: client.conn,
			}
			val = starlark.NewBuiltin("", call.CallInternal)
		}
		client.attrs[key] = val
	}

	return client, nil
}

type unaryCall struct {
	name string
	md   protoreflect.MethodDescriptor
	conn *grpc.ClientConn
}

// String implements part of the starlark.Value interface
func (c *unaryCall) String() string { return fmt.Sprintf("unary-rpc <%s>", c.name) }

// Type implements part of the starlark.Value interface
func (*unaryCall) Type() string { return "grpc.UnaryClientCall" }

// Freeze implements part of the starlark.Value interface
func (*unaryCall) Freeze() {} // immutable

// Truth implements part of the starlark.Value interface
func (*unaryCall) Truth() starlark.Bool { return starlark.False }

// Hash implements part of the starlark.Value interface
func (c *unaryCall) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

func (c *unaryCall) CallInternal(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%s requires a single argument (the request proto)", c)
	}
	request, ok := protomodule.AsProtoMessage(args.Index(0))
	if !ok {
		return nil, fmt.Errorf("failed to convert request argument to ProtoMessage: %v", args.Index(0))
	}
	ctx := context.Background()
	response := dynamicpb.NewMessage(c.md.Output())
	if err := c.conn.Invoke(ctx, c.name, request, response); err != nil {
		return nil, err
	}
	out, err := protomodule.NewMessage(response)
	if err != nil {
		return nil, err
	}
	return out, nil
}
