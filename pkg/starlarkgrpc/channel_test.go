package starlarkgrpc

import (
	"testing"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
)

func TestStarlarkGrpcChannelExpr(t *testing.T) {
	moduletest.ExprTests(t, testModule, []*moduletest.ExprTest{
		{
			Expr: "grpc.Channel",
			Want: `<built-in function grpc.Channel>`,
		},
		{
			Expr:    "grpc.Channel()",
			WantErr: `grpc.Channel: missing argument for target`,
		},
	})
}
