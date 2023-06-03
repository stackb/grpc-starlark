package thread

import (
	"fmt"
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
		"sleep":    starlark.NewBuiltin("thread.sleep", sleep),
		"cancel":   starlark.NewBuiltin("thread.cancel", cancel),
		"timeout":  starlark.NewBuiltin("thread.timeout", timeout),
		"interval": starlark.NewBuiltin("thread.interval", interval),
		"name":     starlark.NewBuiltin("thread.name", name),
	},
}

func sleep(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var millis int64
	if err := starlark.UnpackArgs("thread.sleep", args, kwargs,
		"millis", &millis,
	); err != nil {
		return nil, err
	}
	time.Sleep(time.Millisecond * time.Duration(millis))
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

func timeout(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var millis int64
	var fn starlark.Callable
	if err := starlark.UnpackArgs("thread.timeout", args, kwargs,
		"fn", &fn,
		"millis?", &millis,
	); err != nil {
		return nil, err
	}

	var cancelled bool

	go func() {
		time.Sleep(time.Millisecond * time.Duration(millis))
		if cancelled {
			return
		}
		thread2 := newThread(thread, fmt.Sprintf("thread.timeout(%d)", millis))
		_, err := starlark.Call(thread2, fn, starlark.Tuple{}, []starlark.Tuple{})
		if err != nil && !cancelled {
			thread2.Cancel(fmt.Sprintf("error invoking %s callback function: %v", b.Name(), err))
		}
	}()

	return starlark.NewBuiltin("thread.timeout.cancel", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		cancelled = true
		return starlark.None, nil
	}), nil
}

func interval(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var millis int64
	var fn starlark.Callable
	if err := starlark.UnpackArgs("thread.interval", args, kwargs,
		"fn", &fn,
		"millis?", &millis,
	); err != nil {
		return nil, err
	}

	var cancelled bool

	go func() {
		for {
			time.Sleep(time.Millisecond * time.Duration(millis))
			if cancelled {
				return
			}
			thread2 := newThread(thread, fmt.Sprintf("thread.interval(%d)", millis))
			_, err := starlark.Call(thread2, fn, starlark.Tuple{}, []starlark.Tuple{})
			if err != nil && !cancelled {
				thread2.Cancel(fmt.Sprintf("error invoking %s callback function: %v", b.Name(), err))
			}
		}
	}()

	return starlark.NewBuiltin("thread.interval.cancel", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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
