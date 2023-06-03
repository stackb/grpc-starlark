package starlarkthread

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/starlark"
)

func TestThreadModule(t *testing.T) {
	testCases := map[string]struct {
		input       string
		wantElapsed time.Duration
		wantErr     string
		wantPrinted string
		want        string
	}{
		// thread.sleep
		"sleep function": {
			input: "thread.sleep",
			want:  "<built-in function thread.sleep>",
		},
		"invoke sleep without args": {
			input:   "thread.sleep()",
			wantErr: "thread.sleep: missing argument for millis",
		},
		"invoke sleep with negative number": {
			input: "thread.sleep(-1)",
			want:  "None",
		},
		"invoke sleep with zero": {
			input: "thread.sleep(0)",
			want:  "None",
		},
		"invoke sleep with kwarg": {
			input: "thread.sleep(millis = 0)",
			want:  "None",
		},
		"invoke sleep with string": {
			input:   "thread.sleep('foo')",
			wantErr: "thread.sleep: for parameter millis: got string, want int",
		},
		"invoke sleep, assert delay": {
			input:       "thread.sleep(1000)",
			wantElapsed: time.Millisecond * 1000,
			want:        "None",
		},
		// thread.cancel
		"cancel function": {
			input: "thread.cancel",
			want:  "<built-in function thread.cancel>",
		},
		"invoke cancel without reason": {
			input:   "thread.cancel()",
			wantErr: "Starlark computation cancelled: ",
		},
		"invoke cancel with reason": {
			input:   "thread.cancel(reason = 'testing')",
			wantErr: "Starlark computation cancelled: testing",
		},
		// thread.timeout
		"timeout function": {
			input: "thread.timeout",
			want:  "<built-in function thread.timeout>",
		},
		"invoke timeout to demonstrate return value is the cancel function": {
			input: "thread.timeout(lambda: print('ok'))",
			want:  "<built-in function thread.timeout.cancel>",
		},
		"invoke timeout without delay": {
			input:       "thread.timeout(lambda: print('ok')) and thread.sleep(100)",
			wantPrinted: "ok\n",
			want:        "None",
		},
		"invoke timeout with cancellation": {
			input: "thread.timeout(lambda: print('never'))()",
			want:  "None",
		},
		"invoke timeout with delay": {
			input:       "thread.timeout(fn = lambda: print('ok'), millis = 100) and thread.sleep(200)",
			wantPrinted: "ok\n",
			want:        "None",
		},
		// thread.interval
		"interval function": {
			input: "thread.interval",
			want:  "<built-in function thread.interval>",
		},
		"invoke interval to demonstrate return value is the cancel function": {
			input: "thread.interval(lambda: print('ok'))",
			want:  "<built-in function thread.interval.cancel>",
		},
		"invoke interval with delay": {
			input:       "thread.interval(lambda: print('ok'), millis = 40) and thread.sleep(100)",
			wantPrinted: "ok\nok\n",
			want:        "None",
		},
		"invoke interval with cancellation": {
			input: "thread.interval(lambda: print('never'))()",
			want:  "None",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			start := time.Now()
			thread := new(starlark.Thread)
			var gotPrinted bytes.Buffer
			thread.Print = func(thread *starlark.Thread, msg string) {
				gotPrinted.WriteString(msg)
				gotPrinted.WriteString("\n")
			}
			value, err := starlark.Eval(
				thread,
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
					t.Fatalf("error (-want +got):\n%s", diff)
				}
				return
			}

			gotElapsed := time.Since(start)
			got := value.String()

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("expr (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantPrinted, gotPrinted.String()); diff != "" {
				t.Errorf("print (-want +got):\n%s", diff)
			}

			if gotElapsed < tc.wantElapsed {
				t.Errorf("expected test case time elapsed to be at least %v (got %v)", tc.wantElapsed, gotElapsed)
			}
		})
	}
}
