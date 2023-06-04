package starlarkgrpc

import (
	"testing"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
)

func TestStarlarkGrpcStatusExpr(t *testing.T) {
	moduletest.ExprTests(t, testModule, []*moduletest.ExprTest{
		{
			Expr: "grpc.status",
			Want: "grpc.status(ABORTED = 10, ALREADYEXISTS = 6, CANCELED = 1, DATALOSS = 15, DEADLINEEXCEEDED = 4, FAILEDPRECONDITION = 9, INTERNAL = 13, INVALIDARGUMENT = 3, NOTFOUND = 5, OK = 0, OUTOFRANGE = 11, PERMISSIONDENIED = 7, RESOURCEEXHAUSTED = 8, UNAUTHENTICATED = 16, UNAVAILABLE = 14, UNIMPLEMENTED = 12, UNKNOWN = 2)",
		},
		{
			Expr: "grpc.status.OK",
			Want: "0",
		},
		{
			Expr:    "grpc.status.Foo",
			WantErr: "grpc.status struct has no .Foo attribute",
		},
	})
}
