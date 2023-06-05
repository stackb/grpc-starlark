package starlarkthread

import (
	"fmt"
	"time"

	libtime "go.starlark.net/lib/time"
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
		"sleep":  starlark.NewBuiltin("thread.sleep", sleep),
		"cancel": starlark.NewBuiltin("thread.cancel", cancel),
		"defer":  starlark.NewBuiltin("thread.defer", deferFunc),
		"name":   starlark.NewBuiltin("thread.name", name),
	},
}

func sleep(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var duration libtime.Duration
	if err := starlark.UnpackArgs("thread.sleep", args, kwargs,
		"duration", &duration,
	); err != nil {
		return nil, err
	}
	time.Sleep(time.Duration(duration))
	return starlark.None, nil
}

func cancel(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var reason string
	if err := starlark.UnpackArgs("thread.cancel", args, kwargs,
		"reason?", &reason,
	); err != nil {
		return nil, err
	}
	thread.Cancel(reason)
	return starlark.None, nil
}

func deferFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var duration libtime.Duration
	var count int64 = 1
	var fn starlark.Callable
	if err := starlark.UnpackArgs("thread.defer", args, kwargs,
		"fn", &fn,
		"delay?", &duration,
		"count?", &count,
	); err != nil {
		return nil, err
	}

	var cancelled bool

	go func() {
		for count > 0 {
			if cancelled {
				return
			}
			time.Sleep(time.Duration(duration))
			if cancelled {
				return
			}
			thread2 := newThread(thread, fmt.Sprintf("thread.defer(%d)", duration))
			_, err := starlark.Call(thread2, fn, starlark.Tuple{}, []starlark.Tuple{})
			if err != nil && !cancelled {
				thread2.Cancel(fmt.Sprintf("error invoking %s callback function: %v", b.Name(), err))
			}
			count--
		}
	}()

	return starlark.NewBuiltin("thread.defer.cancel", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		cancelled = true
		return starlark.None, nil
	}), nil
}

func name(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.String(thread.Name), nil
}

func newThread(thread *starlark.Thread, name string) *starlark.Thread {
	newThread := new(starlark.Thread)
	newThread.Name = fmt.Sprintf("%s-%s", thread.Name, name)
	newThread.Print = thread.Print
	newThread.OnMaxSteps = thread.OnMaxSteps
	newThread.Load = thread.Load
	return newThread
}
