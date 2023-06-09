package moduletest

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/starlark"
)

type ExprTestOption func(*ExprTest)

func ExprTests(t *testing.T, globals starlark.StringDict, tt []*ExprTest, options ...ExprTestOption) {
	for _, tc := range tt {
		for _, opt := range options {
			opt(tc)
		}
		t.Run(tc.Expr, func(t *testing.T) {
			tc.Run(t, globals)
		})
	}
}

type ExprTest struct {
	Expr        string            // Input expression
	Env         map[string]string // Optional env vars
	WantErr     string            // Optional expected error
	WantElapsed time.Duration     // Optional expected min test time
	WantPrinted string            // Optional output of 'print'
	Want        string
	Before      func(t *testing.T, globals starlark.StringDict) starlark.StringDict // Optional function to run at the end
	After       func(t *testing.T, value starlark.Value)                            // Optional function to run at the end
}

func (tc *ExprTest) Run(t *testing.T, globals starlark.StringDict) {
	for k, v := range tc.Env {
		os.Setenv(k, v)
	}

	start := time.Now()

	thread := new(starlark.Thread)

	var gotPrinted bytes.Buffer
	thread.Print = func(thread *starlark.Thread, msg string) {
		gotPrinted.WriteString(msg)
		gotPrinted.WriteString("\n")
	}

	if tc.Before != nil {
		globals = tc.Before(t, globals)
	}

	value, err := starlark.Eval(
		thread,
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
	gotElapsed := time.Since(start)

	if diff := cmp.Diff(tc.Want, got); diff != "" {
		t.Errorf("expr (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(strings.TrimSpace(tc.WantPrinted), strings.TrimSpace(gotPrinted.String())); diff != "" {
		t.Errorf("print (-want +got):\n%s", diff)
	}

	if gotElapsed < tc.WantElapsed {
		t.Errorf("expected test case time elapsed to be at least %v (got %v)", tc.WantElapsed, gotElapsed)
	}

	if tc.After != nil {
		tc.After(t, value)
	}
}
