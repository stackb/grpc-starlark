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
		{
			input: "net.Listener(network = 'tcp', address='localhost:1301')",
			want:  "<net.Listener tcp 127.0.0.1:1301>",
		},
		{
			input: "net.Listener().network",
			want:  `"tcp"`,
		},
		{
			input: "net.Listener(address='localhost:1300').address",
			want:  `"127.0.0.1:1300"`,
		},
		{
			input: "net.Listener(address=':1302')",
			want:  "<net.Listener tcp [::]:1302>",
		},
		{
			input: "net.Listener(address=':1303')",
			want:  "<net.Listener tcp [::]:1303>",
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
