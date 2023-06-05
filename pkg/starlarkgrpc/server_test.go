package starlarkgrpc

import (
	"testing"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
)

func TestStarlarkGrpcServerExpr(t *testing.T) {
	moduletest.ExprTests(t, testModule, []*moduletest.ExprTest{
		// grpc.Server
		{
			Expr: "grpc.Server",
			Want: `<built-in function grpc.Server>`,
		},
		{
			Expr: "grpc.Server()",
			Want: `<grpc.Server []>`,
		},
		// grpc.Server.start
		{
			Expr: "grpc.Server().start",
			Want: `<built-in function grpc.Server.start>`,
		},
		{
			Expr:    "grpc.Server().start()",
			WantErr: `grpc.Server.start: got 0 arguments, want 1`,
		},
		// grpc.Server.stop
		{
			Expr: "grpc.Server().stop",
			Want: `<built-in function grpc.Server.stop>`,
		},
		{
			Expr: "grpc.Server().stop()",
			Want: `None`,
		},
		{
			Expr: "grpc.Server().stop(graceful = False)",
			Want: `None`,
		},
		// grpc.Server.register
		{
			Expr: "grpc.Server().register",
			Want: `<built-in function grpc.Server.register>`,
		},
		{
			Expr:    "grpc.Server().register()",
			WantErr: `grpc.Server.register: missing argument for service`,
		},
		{
			Expr:    "grpc.Server().register('example.routeguide.Routeguide', {})",
			WantErr: `unknown service: example.routeguide.Routeguide (known: [])`,
		},
		{
			Expr:    "grpc.Server().register(service = 'example.routeguide.Routeguide', handlers = {})",
			WantErr: `unknown service: example.routeguide.Routeguide (known: [])`,
		},
	})
}
