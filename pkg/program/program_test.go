package program

import (
	"bytes"
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stackb/grpc-starlark/pkg/protodescriptorset"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protodesc"
)

//go:embed routeguide_proto_descriptor.pb
var routeguideProtodesciptorSet []byte

func TestProgram(t *testing.T) {
	testCases := []struct {
		program string
		env     map[string]string
		wantErr string
		want    string
	}{
		{
			program: "print(grpc.status.OK)",
			want:    "0",
		},
	}

	for _, tc := range testCases {
		var gotPrinted bytes.Buffer
		thread := new(starlark.Thread)
		thread.Print = func(thread *starlark.Thread, msg string) {
			gotPrinted.WriteString(msg)
			gotPrinted.WriteString("\n")
		}
		for k, v := range tc.env {
			os.Setenv(k, v)
		}

		pds, err := protodescriptorset.Parse(routeguideProtodesciptorSet)
		if err != nil {
			t.Fatal(err)
		}
		files, err := protodesc.NewFiles(pds)
		if err != nil {
			t.Fatal(err)
		}
		globals := NewPredeclared(files)

		_, err = starlark.ExecFile(
			thread,
			"<in-memory>",
			strings.NewReader(tc.program),
			globals,
		)

		if err != nil {
			if tc.wantErr == "" {
				t.Error("unexpected error: ", err)
				continue
			}
			gotErr := err.Error()
			if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
			continue
		}

		got := strings.TrimSpace(gotPrinted.String())
		want := strings.TrimSpace(tc.want)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
