package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/stackb/grpc-starlark/pkg/starlarkgrpc"
)

type Server struct {
	server   *grpc.Server
	methods  map[string]protoreflect.MethodDescriptor
	handlers starlarkgrpc.HandlerMap
}

func New(files *protoregistry.Files) (*Server, error) {
	s := &Server{
		handlers: make(starlarkgrpc.HandlerMap),
		methods:  make(map[string]protoreflect.MethodDescriptor),
	}
	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.InterceptUnary))

	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			sd := services.Get(i)
			methods := sd.Methods()
			gsd := grpc.ServiceDesc{ServiceName: string(sd.FullName()), HandlerType: (*interface{})(nil)}
			for j := 0; j < methods.Len(); j++ {
				method := methods.Get(j)
				name := fmt.Sprintf("/%s/%s", sd.FullName(), method.Name())
				if method.IsStreamingServer() && method.IsStreamingClient() {
					gsd.Streams = append(gsd.Streams, grpc.StreamDesc{
						StreamName:    string(method.Name()),
						ServerStreams: true,
						ClientStreams: true,
						Handler:       s.HandleStream,
					})
					log.Printf("Registered %s (bidi stream):", name)
				} else if method.IsStreamingServer() {
					gsd.Streams = append(gsd.Streams, grpc.StreamDesc{
						StreamName:    string(method.Name()),
						ServerStreams: true,
						Handler:       s.HandleStream,
					})
					log.Printf("Registered %s (server stream):", name)
				} else if method.IsStreamingClient() {
					gsd.Streams = append(gsd.Streams, grpc.StreamDesc{
						StreamName:    string(method.Name()),
						ClientStreams: true,
						Handler:       s.HandleStream,
					})
					log.Printf("Registered %s (client stream):", name)
				} else {
					gsd.Methods = append(gsd.Methods, grpc.MethodDesc{
						MethodName: string(method.Name()),
						Handler:    s.HandleMethod,
					})
					log.Printf("Registered %s (unary method):", name)
				}
				s.methods[name] = method
			}

			s.server.RegisterService(&gsd, s)
		}
		return true
	})

	return s, nil
}

func (s *Server) OnHandler(handler *starlarkgrpc.Handler) error {
	if _, ok := s.methods[handler.Name()]; ok {
		log.Printf("Registered handler for %s", handler.Name())
		s.handlers[handler.Name()] = handler
		return nil
	}
	return fmt.Errorf("error: starlark handler %q has no matching gRPC method", handler.Name())
}

func (s *Server) Start(l net.Listener) error {
	go s.server.Serve(l)
	return nil
}

func (s *Server) InterceptUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("InterceptUnary %s", info.FullMethod)
	resp, err = handler(ctx, req)
	return
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}

// HandleStream implements grpc.StreamHandler for handling of server-streaming
// calls.
func (s *Server) HandleStream(srv interface{}, ss grpc.ServerStream) error {
	stream := grpc.ServerTransportStreamFromContext(ss.Context())

	method, ok := s.methods[stream.Method()]
	if !ok {
		log.Println("method not found:", stream.Method())
		return status.Error(codes.Unimplemented, stream.Method())
	}

	handler, ok := s.handlers[stream.Method()]
	if !ok {
		log.Println("handler not found:", stream.Method())
		return status.Error(codes.Unimplemented, stream.Method())
	}

	var request protoreflect.ProtoMessage
	if method.IsStreamingServer() && !method.IsStreamingClient() {
		request = dynamicpb.NewMessage(method.Input())
		if err := ss.RecvMsg(request); err != nil {
			return err
		}
	}

	response, err := handler.Handle(method, request, ss)
	if err != nil {
		log.Printf("handler return value error: %v", err)
		return err
	}

	log.Println("stream response:", response)

	if method.IsStreamingClient() && !method.IsStreamingServer() {
		if err := ss.SendMsg(response); err != nil {
			return err
		}
	}

	return nil
}

// HandleStream implements grpc.methodHandler for handling of unary calls.
func (s *Server) HandleMethod(srv interface{}, ctx context.Context, decode func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	stream := grpc.ServerTransportStreamFromContext(ctx)
	log.Println("handling:", stream.Method())

	method, ok := s.methods[stream.Method()]
	if !ok {
		log.Println("method not found:", stream.Method())
		return nil, status.Error(codes.Unimplemented, stream.Method())
	}

	handler, ok := s.handlers[stream.Method()]
	if !ok {
		log.Println("handler not found:", stream.Method())
		return nil, status.Error(codes.Unimplemented, stream.Method())
	}

	input := dynamicpb.NewMessage(method.Input())
	if err := decode(input); err != nil {
		return nil, err
	}

	response, err := handler.Handle(method, input, nil)
	if err != nil {
		log.Printf("handler return value error: %v", err)
		return nil, err
	}

	return response, err
}
