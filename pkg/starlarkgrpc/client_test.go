package starlarkgrpc

import (
	"testing"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
)

func TestStarlarkGrpcClientExpr(t *testing.T) {
	moduletest.ExprTests(t, testModule, []*moduletest.ExprTest{
		{
			Expr: "grpc.Client",
			Want: `<built-in function grpc.Client>`,
		},
		{
			Expr:    "grpc.Client()",
			WantErr: `grpc.Client: missing argument for service`,
		},
		{
			Expr:    "grpc.Client(service = 'example.routeguide.Routeguide')",
			WantErr: `grpc.Client: missing argument for channel`,
		},
		{
			Expr:    "grpc.Client(service = 'example.routeguide.Routeguide', channel = grpc.Channel(':0'))",
			WantErr: `unknown service: example.routeguide.Routeguide (known: [])`,
		},
	})
}
