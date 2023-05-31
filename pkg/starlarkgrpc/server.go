package starlarkgrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"go.starlark.net/starlark"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

// grpcServer implements starlark.Value for a grpc.Server.
type grpcServer struct {
	server   *grpc.Server
	handlers MethodHandlerMap
}

// String implements part of the starlark.Value interface
func (*grpcServer) String() string { return "GrpcServer" }

// Type implements part of the starlark.Value interface
func (*grpcServer) Type() string { return "GrpcServer" }

// Freeze implements part of the starlark.Value interface
func (*grpcServer) Freeze() {} // immutable

// Truth implements part of the starlark.Value interface
func (*grpcServer) Truth() starlark.Bool { return starlark.False }

// Hash implements part of the starlark.Value interface
func (c *grpcServer) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

// AttrNames implements part of the starlark.HasAttrs interface
func (c *grpcServer) AttrNames() []string {
	return []string{"start", "stop", "register"}
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
	if args.Len() != 1 {
		return nil, fmt.Errorf("grpc.Server.start requires exactly one argument (a listener)")
	}
	value := args.Index(0)
	lis, ok := value.(net.Listener)
	if !ok {
		return nil, fmt.Errorf("grpc.Server.start argument must implement net.Listener (got %T)", value)
	}
	go func() {
		log.Printf("grpc.Server listening on %v", lis.Addr())
		if err := c.server.Serve(lis); err != nil {
			log.Fatalln("grpc.Server error:", err)
		}
	}()
	return starlark.None, nil
}

func (c *grpcServer) stop(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	c.server.GracefulStop()
	return starlark.None, nil
}

func (c *grpcServer) register(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var serviceName string
	var mappings *starlark.Dict
	if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 2, &serviceName, &mappings); err != nil {
		return nil, err
	}

	var sd protoreflect.ServiceDescriptor
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			s := services.Get(i)
			if string(s.FullName()) == serviceName {
				sd = s
				return false
			}
		}
		return true
	})
	if sd == nil {
		return nil, fmt.Errorf("grpc.register error: unknown service: %s (known: %v)", serviceName, serviceNames(protoregistry.GlobalFiles))
	}
	gsd := grpc.ServiceDesc{ServiceName: string(sd.FullName()), HandlerType: (*interface{})(nil)}

	methods := sd.Methods()
	for _, key := range mappings.Keys() {
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
			method:   method,
		}

		key := fmt.Sprintf("/%s/%s", sd.FullName(), method.Name())
		if method.IsStreamingServer() && method.IsStreamingClient() {
			gsd.Streams = append(gsd.Streams, grpc.StreamDesc{
				StreamName:    string(method.Name()),
				ServerStreams: true,
				ClientStreams: true,
				Handler:       c.HandleStream,
			})
			log.Printf("grpc.Server: Registered %s (bidi stream):", key)
		} else if method.IsStreamingServer() {
			gsd.Streams = append(gsd.Streams, grpc.StreamDesc{
				StreamName:    string(method.Name()),
				ServerStreams: true,
				Handler:       c.HandleStream,
			})
			log.Printf("grpc.Server: Registered %s (server stream):", key)
		} else if method.IsStreamingClient() {
			gsd.Streams = append(gsd.Streams, grpc.StreamDesc{
				StreamName:    string(method.Name()),
				ClientStreams: true,
				Handler:       c.HandleStream,
			})
			log.Printf("grpc.Server: Registered %s (client stream):", key)
		} else {
			gsd.Methods = append(gsd.Methods, grpc.MethodDesc{
				MethodName: string(method.Name()),
				Handler:    c.HandleMethod,
			})
			log.Printf("grpc.Server: Registered %s (unary method):", key)
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
	if handler.method.IsStreamingServer() && !handler.method.IsStreamingClient() {
		request = dynamicpb.NewMessage(handler.method.Input())
		if err := ss.RecvMsg(request); err != nil {
			return err
		}
	}

	response, err := handler.Handle(handler.method, request, ss)
	if err != nil {
		log.Printf("handler return value error: %v", err)
		return err
	}

	log.Println("stream response:", response)

	if handler.method.IsStreamingClient() && !handler.method.IsStreamingServer() {
		if err := ss.SendMsg(response); err != nil {
			return err
		}
	}

	return nil
}

// HandleStream implements grpc.methodHandler for handling of unary calls.
func (s *grpcServer) HandleMethod(srv interface{}, ctx context.Context, decode func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	stream := grpc.ServerTransportStreamFromContext(ctx)
	log.Println("grpc.Server handle method:", stream.Method())

	handler, ok := s.handlers[stream.Method()]
	if !ok {
		log.Println("handler not found:", stream.Method())
		return nil, status.Error(codes.Unimplemented, stream.Method())
	}

	input := dynamicpb.NewMessage(handler.method.Input())
	if err := decode(input); err != nil {
		return nil, err
	}

	response, err := handler.Handle(handler.method, input, nil)
	if err != nil {
		log.Printf("handler return value error: %v", err)
		return nil, err
	}

	return response, err
}

func newServerFunction() goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

		if err := starlark.UnpackArgs("Server", args, kwargs); err != nil {
			return nil, err
		}

		value := &grpcServer{
			server:   grpc.NewServer(),
			handlers: make(map[string]*MethodHandler),
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
