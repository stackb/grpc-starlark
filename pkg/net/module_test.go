package net

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/starlark"
)

func TestNetModule(t *testing.T) {
	testCases := []struct {
		input   string
		wantErr string
		want    string
	}{
		{
			input: "net.Listener",
			want:  "<built-in function net.Listener>",
		},
	}

	for _, tc := range testCases {
		value, err := starlark.Eval(
			new(starlark.Thread),
			"<expr>",
			tc.input,
			starlark.StringDict{
				"net": Module,
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
