package starlarknet

import (
	"testing"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
	"go.starlark.net/starlark"
)

func TestNetModule(t *testing.T) {
	moduletest.ExprTests(t, starlark.StringDict{
		"net": Module,
	}, []*moduletest.ExprTest{
		{
			Expr: "net.Listener",
			Want: "<built-in function net.Listener>",
		},
		{
			Expr: "net.Listener(network = 'tcp', address='localhost:1301')",
			Want: "<net.Listener tcp 127.0.0.1:1301>",
		},
		{
			Expr: "net.Listener().network",
			Want: `"tcp"`,
		},
		{
			Expr: "net.Listener(address='localhost:1300').address",
			Want: `"127.0.0.1:1300"`,
		},
		{
			Expr: "net.Listener(address=':1302')",
			Want: "<net.Listener tcp [::]:1302>",
		},
		{
			Expr: "net.Listener(address=':1303')",
			Want: "<net.Listener tcp [::]:1303>",
		},
	})
}
