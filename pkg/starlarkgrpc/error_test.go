package starlarkgrpc

import (
	"testing"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
)

func TestStarlarkGrpcErrorExpr(t *testing.T) {
	moduletest.ExprTests(t, testModule, []*moduletest.ExprTest{
		{
			Expr: "grpc.Error",
			Want: "<built-in function grpc.Error>",
		},
		{
			Expr: "grpc.Error()",
			Want: `grpc.Error(code = 2, message = "")`,
		},
		{
			Expr: "grpc.Error(code = grpc.status.ABORTED, message = 'user')",
			Want: `grpc.Error(code = 10, message = "user")`,
		},
	})
}
