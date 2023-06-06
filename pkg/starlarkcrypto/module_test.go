package starlarkcrypto

import (
	"testing"

	"go.starlark.net/starlark"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
)

func TestTlsModule(t *testing.T) {
	moduletest.ExprTests(t, starlark.StringDict{
		"crypto": Module,
	}, []*moduletest.ExprTest{
		{
			Expr: "crypto.tls.Certificate",
			Want: "<built-in function tls.Certificate>",
		},
	})
}
