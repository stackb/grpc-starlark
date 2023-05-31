package thread

import (
	"time"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module thread is a Starlark module of thread-related functions and constants. The
// module defines the following functions:
//
//	sleep(millis) - Sleeps the current thread for the given number of milliseconds
var Module = &starlarkstruct.Module{
	Name: "thread",
	Members: starlark.StringDict{
		"sleep": starlark.NewBuiltin("sleep", sleep),
	},
}

func sleep(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var millis int64
	if err := starlark.UnpackArgs("sleep", args, kwargs,
		"millis", &millis,
	); err != nil {
		return nil, err
	}
	time.Sleep(time.Millisecond * time.Duration(millis))
	return starlark.None, nil
}
