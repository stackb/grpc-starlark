package starlarkgrpc

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/starlark"
)

func TestStarlarkGrpcModule(t *testing.T) {
	testCases := []struct {
		input   string
		env     map[string]string
		wantErr string
		want    string
	}{
		{
			input: "grpc.Error",
			want:  "<built-in function grpc.Error>",
		},
	}

	for _, tc := range testCases {
		for k, v := range tc.env {
			os.Setenv(k, v)
		}
		value, err := starlark.Eval(
			new(starlark.Thread),
			"<expr>",
			tc.input,
			starlark.StringDict{
				"grpc": Module,
			},
		)
		if err != nil {
			if tc.wantErr == "" {
				t.Fatal("unexpected error: ", err)
			}
			gotErr := err.Error()
			if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
				t.Fatalf("(-want +got):\n%s", diff)
			}
			return
		}

		got := value.String()

		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
