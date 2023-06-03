package starlarkos

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/starlark"
)

func TestOsModule(t *testing.T) {
	testCases := []struct {
		input   string
		env     map[string]string
		wantErr string
		want    string
	}{
		{
			input: "os.getenv",
			want:  "<built-in function os.getenv>",
		},
		{
			input:   "os.getenv()",
			wantErr: "os.getenv: missing argument for name",
		},
		{
			input: "os.getenv('FOO')",
			want:  "None",
		},
		{
			input: "os.getenv('FOO')",
			env:   map[string]string{"FOO": "BAR"},
			want:  "BAR",
		},
		{
			input: "os.getenv('FOO')",
			env:   map[string]string{"FOO": ""},
			want:  "",
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
				"os": Module,
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
