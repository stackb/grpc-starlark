package starlarkgrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"go.starlark.net/starlark"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

// grpcServer implements starlark.Value for a grpc.Server.
type grpcServer struct {
	server   *grpc.Server
	files    *protoregistry.Files
	handlers HandlerMap
}

// String implements part of the starlark.Value interface
func (s *grpcServer) String() string {
	return fmt.Sprintf("<grpc.Server %v>", grpcServerServiceInfoNames(s.server))
}

// Type implements part of the starlark.Value interface
func (*grpcServer) Type() string { return "grpc.Server" }

// Freeze implements part of the starlark.Value interface
func (*grpcServer) Freeze() {} // immutable

// Truth implements part of the starlark.Value interface
func (*grpcServer) Truth() starlark.Bool { return starlark.True }

// Hash implements part of the starlark.Value interface
func (c *grpcServer) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

// AttrNames implements part of the starlark.HasAttrs interface
func (c *grpcServer) AttrNames() []string {
	return []string{"register", "start", "stop"}
}

// Attr implements part of the starlark.HasAttrs interface
func (c *grpcServer) Attr(name string) (starlark.Value, error) {
	switch name {
	case "start":
		return starlark.NewBuiltin("grpc.Server.start", c.start), nil
	case "stop":
		return starlark.NewBuiltin("grpc.Server.stop", c.stop), nil
	case "register":
		return starlark.NewBuiltin("grpc.Server.register", c.register), nil
	}
	return nil, nil
}

func (c *grpcServer) start(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var listener net.Listener
	if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &listener); err != nil {
		return nil, err
	}
	if err := c.server.Serve(listener); err != nil {
		log.Fatalln("grpc.Server error:", err)
	}
	return starlark.None, nil
}

func (c *grpcServer) stop(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	graceful := true
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
		"graceful?", &graceful,
	); err != nil {
		return nil, err
	}
	if graceful {
		c.server.GracefulStop()
	} else {
		c.server.Stop()
	}
	return starlark.None, nil
}

func (c *grpcServer) register(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var serviceName string
	var handlers *starlark.Dict
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
		"service", &serviceName,
		"handlers", &handlers,
	); err != nil {
		return nil, err
	}

	sd, ok := serviceDescriptorByName(c.files, serviceName)
	if !ok {
		return nil, fmt.Errorf("unknown service: %s (known: %v)", serviceName, serviceNames(c.files))
	}

	gsd := grpc.ServiceDesc{ServiceName: string(sd.FullName()), HandlerType: (*interface{})(nil)}

	methods := sd.Methods()
	for _, key := range handlers.Keys() {
		name, ok := key.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("%s: register error: dict key should be a fully-qualified method name (got %T)", fn.Name(), key)
		}

		var method protoreflect.MethodDescriptor
		for i := 0; i < methods.Len(); i++ {
			md := methods.Get(i)
			if string(md.Name()) == name.GoString() {
				method = md
				break
			}
		}
		if method == nil {
			return nil, fmt.Errorf("grpc.register error: unknown method %s for service %s", name.GoString(), serviceName)
		}

		value, ok, err := handlers.Get(key)
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

		handler := &Handler{
			name:     name.GoString(),
			fn:       callable,
			reporter: thread.Print,
			md:       method,
		}

		key := fmt.Sprintf("/%s/%s", sd.FullName(), method.Name())
		if method.IsStreamingServer() && method.IsStreamingClient() {
			gsd.Streams = append(gsd.Streams, grpc.StreamDesc{
				StreamName:    string(method.Name()),
				ServerStreams: true,
				ClientStreams: true,
				Handler:       c.HandleStream,
			})
		} else if method.IsStreamingServer() {
			gsd.Streams = append(gsd.Streams, grpc.StreamDesc{
				StreamName:    string(method.Name()),
				ServerStreams: true,
				Handler:       c.HandleStream,
			})
		} else if method.IsStreamingClient() {
			gsd.Streams = append(gsd.Streams, grpc.StreamDesc{
				StreamName:    string(method.Name()),
				ClientStreams: true,
				Handler:       c.HandleStream,
			})
		} else {
			gsd.Methods = append(gsd.Methods, grpc.MethodDesc{
				MethodName: string(method.Name()),
				Handler:    c.HandleMethod,
			})
		}
		c.handlers[key] = handler
	}

	c.server.RegisterService(&gsd, c)

	return starlark.None, nil
}

// HandleStream implements grpc.StreamHandler for handling of server-streaming
// calls.
func (s *grpcServer) HandleStream(srv interface{}, ss grpc.ServerStream) error {
	stream := grpc.ServerTransportStreamFromContext(ss.Context())

	handler, ok := s.handlers[stream.Method()]
	if !ok {
		log.Println("handler not found:", stream.Method())
		return status.Error(codes.Unimplemented, stream.Method())
	}

	var request protoreflect.ProtoMessage
	if handler.md.IsStreamingServer() && !handler.md.IsStreamingClient() {
		request = dynamicpb.NewMessage(handler.md.Input())
		if err := ss.RecvMsg(request); err != nil {
			return err
		}
	}

	response, err := handler.handle(handler.md, request, ss.Context(), ss, nil)
	if err != nil {
		log.Printf("handler return value error: %v", err)
		return err
	}

	if handler.md.IsStreamingClient() && !handler.md.IsStreamingServer() {
		if err := ss.SendMsg(response); err != nil {
			return err
		}
	}

	return nil
}

// HandleStream implements grpc.methodHandler for handling of unary calls.
func (s *grpcServer) HandleMethod(srv interface{}, ctx context.Context, decode func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	sts := grpc.ServerTransportStreamFromContext(ctx)

	handler, ok := s.handlers[sts.Method()]
	if !ok {
		log.Println("handler not found:", sts.Method())
		return nil, status.Error(codes.Unimplemented, sts.Method())
	}

	input := dynamicpb.NewMessage(handler.md.Input())
	if err := decode(input); err != nil {
		return nil, err
	}

	response, err := handler.handle(handler.md, input, ctx, nil, sts)
	if err != nil {
		return nil, err
	}

	return response, err
}

func newServer(files *protoregistry.Files) func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var credentials credentials.TransportCredentials

		if err := starlark.UnpackArgs("grpc.Server", args, kwargs,
			"credentials?", &credentials,
		); err != nil {
			return nil, err
		}
		var options []grpc.ServerOption
		if credentials != nil {
			options = append(options, grpc.Creds(credentials))
		}
		value := &grpcServer{
			files:    files,
			server:   grpc.NewServer(options...),
			handlers: make(map[string]*Handler),
		}
		return value, nil
	}
}

func serviceNames(files *protoregistry.Files) (names []string) {
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			sd := services.Get(i)
			names = append(names, string(sd.FullName()))
		}
		return true
	})
	return
}

func grpcServerServiceInfoNames(server *grpc.Server) (names []string) {
	for name := range server.GetServiceInfo() {
		names = append(names, name)
	}
	return
}

func serviceDescriptorByName(files *protoregistry.Files, name string) (protoreflect.ServiceDescriptor, bool) {
	var sd protoreflect.ServiceDescriptor
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			s := services.Get(i)
			if string(s.FullName()) == name {
				sd = s
				return false
			}
		}
		return true
	})
	if sd == nil {
		return nil, false
	}
	return sd, true
}

func methodDescriptorByName(files *protoregistry.Files, serviceName, methodName string) (md protoreflect.MethodDescriptor, ok bool) {
	service, ok := serviceDescriptorByName(files, serviceName)
	if !ok {
		return nil, false
	}
	methods := service.Methods()
	for i := 0; i < methods.Len(); i++ {
		v := methods.Get(i)
		if string(v.Name()) == string(methodName) {
			return v, true
		}
	}
	return nil, false
}
