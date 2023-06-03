package moduletest

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/starlark"
)

type ExprTests []*ExprTest

func (tt ExprTests) TestAll(t *testing.T, globals starlark.StringDict) {
	for _, tc := range tt {
		t.Run(tc.Expr, func(t *testing.T) {
			tc.Test(t, globals)
		})
	}
}

type ExprTest struct {
	Expr    string
	Env     map[string]string
	WantErr string
	Want    string
}

func (tc *ExprTest) Test(t *testing.T, globals starlark.StringDict) {
	for k, v := range tc.Env {
		os.Setenv(k, v)
	}
	value, err := starlark.Eval(
		new(starlark.Thread),
		"<expr>",
		tc.Expr,
		globals,
	)
	if err != nil {
		if tc.WantErr == "" {
			t.Fatal("unexpected error: ", err)
		}
		gotErr := err.Error()
		if diff := cmp.Diff(tc.WantErr, gotErr); diff != "" {
			t.Fatalf("(-want +got):\n%s", diff)
		}
		return
	}

	got := value.String()
	if diff := cmp.Diff(tc.Want, got); diff != "" {
		t.Errorf("expr (-want +got):\n%s", diff)
	}
}
