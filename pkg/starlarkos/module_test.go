package starlarkos

import (
	"testing"

	"go.starlark.net/starlark"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
)

func TestOsModule(t *testing.T) {
	moduletest.ExprTests(t, starlark.StringDict{
		"os": Module,
	}, []*moduletest.ExprTest{
		{
			Expr: "os.getenv",
			Want: "<built-in function os.getenv>",
		},
		{
			Expr:    "os.getenv()",
			WantErr: "os.getenv: missing argument for name",
		},
		{
			Expr: "os.getenv('FOO')",
			Want: "None",
		},
		{
			Expr: "os.getenv('FOO')",
			Env:  map[string]string{"FOO": "BAR"},
			Want: `"BAR"`,
		},
		{
			Expr: "os.getenv('FOO')",
			Env:  map[string]string{"FOO": ""},
			Want: `""`,
		},
	})
}
