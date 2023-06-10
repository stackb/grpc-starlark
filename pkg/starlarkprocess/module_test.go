package starlarkprocess

import (
	"testing"
	"time"

	"github.com/stackb/grpc-starlark/pkg/moduletest"
	libtime "go.starlark.net/lib/time"
	"go.starlark.net/starlark"
)

func TestThreadModule(t *testing.T) {
	moduletest.ExprTests(t, starlark.StringDict{
		"thread": NewModule(),
		"time":   libtime.Module,
	}, []*moduletest.ExprTest{
		// thread.sleep
		{
			Expr: "thread.sleep",
			Want: "<built-in function thread.sleep>",
		},
		{
			Expr:    "thread.sleep()",
			WantErr: "thread.sleep: missing argument for duration",
		},
		{
			Expr: "thread.sleep(time.millisecond * -1)",
			Want: "None",
		},
		{
			Expr: "thread.sleep(time.millisecond * 0)",
			Want: "None",
		},
		{
			Expr: "thread.sleep(time.millisecond * 0)",
			Want: "None",
		},
		{
			Expr:    "thread.sleep('foo')",
			WantErr: `thread.sleep: for parameter duration: time: invalid duration "foo"`,
		},
		{
			Expr:        "thread.sleep(1000 * time.millisecond)",
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
		// thread.defer
		{
			Expr: "thread.defer",
			Want: "<built-in function thread.defer>",
		},
		{
			Expr: "thread.defer(lambda: print('ok'))",
			Want: "<built-in function thread.defer.cancel>",
		},
		{
			Expr:        "thread.defer(lambda: print('ok')) and thread.sleep(time.millisecond * 100)",
			WantPrinted: "ok",
			Want:        "None",
		},
		{
			Expr: "thread.defer(lambda: print('never'))()",
			Want: "None",
		},
		{
			Expr:        "thread.defer(fn = lambda: print('ok'), delay = time.millisecond * 100) and thread.sleep(time.millisecond * 200)",
			WantPrinted: "ok",
			Want:        "None",
		},
	})
}
