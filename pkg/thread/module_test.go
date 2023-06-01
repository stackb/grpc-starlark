package thread

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/starlark"
)

func TestThreadModule(t *testing.T) {
	testCases := []struct {
		input       string
		wantElapsed time.Duration
		wantErr     string
		want        string
	}{
		{
			input: "thread.sleep",
			want:  "<built-in function thread.sleep>",
		},
		{
<<<<<<< HEAD
			input:   "thread.sleep()",
			wantErr: "thread.sleep: missing argument for millis",
		},
		{
=======
>>>>>>> master
			input: "thread.sleep(-1)",
			want:  "None",
		},
		{
			input: "thread.sleep(0)",
			want:  "None",
		},
		{
			input: "thread.sleep(millis = 0)",
			want:  "None",
		},
		{
			input:   "thread.sleep('foo')",
			wantErr: "thread.sleep: for parameter millis: got string, want int",
		},
		{
			input:       "thread.sleep(1000)",
			wantElapsed: time.Millisecond * 1000,
			want:        "None",
		},
	}

	for _, tc := range testCases {
		start := time.Now()
		value, err := starlark.Eval(
			new(starlark.Thread),
			"<expr>",
			tc.input,
			starlark.StringDict{
				"thread": Module,
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

		gotElapsed := time.Since(start)
		got := value.String()

		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}

		if gotElapsed < tc.wantElapsed {
			t.Errorf("expected test case time elapsed to be at least %v (got %v)", tc.wantElapsed, gotElapsed)
		}
	}
}
