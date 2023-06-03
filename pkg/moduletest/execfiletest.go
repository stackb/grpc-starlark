package moduletest

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	libproto "go.starlark.net/lib/proto"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type ExecFileTests map[string]*ExecFileTest

func (tt ExecFileTests) Run(t *testing.T, files *protoregistry.Files, globals starlark.StringDict) {
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tc.Run(t, name, files, globals)
		})
	}
}

type ExecFileTest struct {
	Expr    string
	Env     map[string]string
	WantErr string
	Want    string
}

func (tc *ExecFileTest) Run(t *testing.T, name string, files *protoregistry.Files, globals starlark.StringDict) {
	var printed bytes.Buffer
	thread := new(starlark.Thread)
	thread.Name = "main:" + name
	thread.Print = func(thread *starlark.Thread, msg string) {
		t.Log(msg)
		printed.WriteString(msg)
		printed.WriteString("\n")
	}

	libproto.SetPool(thread, files)

	_, err := starlark.ExecFile(
		thread,
		"<in-memory>",
		strings.NewReader(tc.Expr),
		globals,
	)

	if err != nil {
		if tc.WantErr == "" {
			t.Error("unexpected error: ", err)
			return
		}
		gotErr := err.Error()
		if diff := cmp.Diff(tc.WantErr, gotErr); diff != "" {
			t.Errorf("error (-want +got):\n%s", diff)
		}
		return
	}

	gotPrinted := strings.TrimSpace(printed.String())
	wantPrinted := strings.TrimSpace(tc.Want)

	if diff := cmp.Diff(wantPrinted, gotPrinted); diff != "" {
		t.Errorf("print (-want +got):\n%s", diff)
	}
}
