package starlarkthread

import (
	"testing"
	"time"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
	"go.starlark.net/starlark"
)

func TestThreadModule(t *testing.T) {
	moduletest.ExprTests(t, starlark.StringDict{
		"thread": Module,
	}, []*moduletest.ExprTest{
		// thread.sleep
		{
			Expr: "thread.sleep",
			Want: "<built-in function thread.sleep>",
		},
		{
			Expr:    "thread.sleep()",
			WantErr: "thread.sleep: missing argument for millis",
		},
		{
			Expr: "thread.sleep(-1)",
			Want: "None",
		},
		{
			Expr: "thread.sleep(0)",
			Want: "None",
		},
		{
			Expr: "thread.sleep(millis = 0)",
			Want: "None",
		},
		{
			Expr:    "thread.sleep('foo')",
			WantErr: "thread.sleep: for parameter millis: got string, want int",
		},
		{
			Expr:        "thread.sleep(1000)",
			WantElapsed: time.Millisecond * 1000,
			Want:        "None",
		},
		// thread.cancel
		{
			Expr: "thread.cancel",
			Want: "<built-in function thread.cancel>",
		},
		{
			Expr:    "thread.cancel()",
			WantErr: "Starlark computation cancelled: ",
		},
		{
			Expr:    "thread.cancel(reason = 'testing')",
			WantErr: "Starlark computation cancelled: testing",
		},
		// thread.timeout
		{
			Expr: "thread.timeout",
			Want: "<built-in function thread.timeout>",
		},
		{
			Expr: "thread.timeout(lambda: print('ok'))",
			Want: "<built-in function thread.timeout.cancel>",
		},
		{
			Expr:        "thread.timeout(lambda: print('ok')) and thread.sleep(100)",
			WantPrinted: "ok",
			Want:        "None",
		},
		{
			Expr: "thread.timeout(lambda: print('never'))()",
			Want: "None",
		},
		{
			Expr:        "thread.timeout(fn = lambda: print('ok'), millis = 100) and thread.sleep(200)",
			WantPrinted: "ok",
			Want:        "None",
		},
		// thread.interval
		{
			Expr: "thread.interval",
			Want: "<built-in function thread.interval>",
		},
		{
			Expr: "thread.interval(lambda: print('ok'))",
			Want: "<built-in function thread.interval.cancel>",
		},
		{
			Expr:        "thread.interval(lambda: print('ok'), millis = 40) and thread.sleep(100)",
			WantPrinted: "ok\nok",
			Want:        "None",
		},
		{
			Expr: "thread.interval(lambda: print('never'))()",
			Want: "None",
		},
	})
}
