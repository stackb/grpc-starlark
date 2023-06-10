package starlarkprocess

import (
	"bytes"
	"os/exec"
	"syscall"

	"github.com/stackb/grpc-starlark/pkg/starlarkutil"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func run(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var command string
	var argv *starlark.List
	var env *starlark.Dict
	var stdin starlark.Bytes
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
		"command", &command,
		"args?", &argv,
		"env?", &env,
		"stdin?", &stdin,
	); err != nil {
		return nil, err
	}

	cmd := exec.Command(command, starlarkListToStringSlice(argv)...)
	if stdin.Len() > 0 {
		cmd.Stdin = bytes.NewBuffer([]byte(stdin))
	}

	var stderr []byte
	var errMsg string
	stdout, err := cmd.Output()
	var exitCode int
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
			stderr = exitError.Stderr
			errMsg = exitError.Error()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH, in
			// this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and
			// format err to string and set to stderr
			exitCode = -1
			errMsg = err.Error()
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	return starlarkstruct.FromStringDict(
		starlarkutil.Symbol(fn.Name()),
		starlark.StringDict{
			"command":   starlark.String(command),
			"args":      args,
			"error":     starlark.String(errMsg),
			"stdout":    starlark.Bytes(stdout),
			"stderr":    starlark.Bytes(stderr),
			"exit_code": starlark.MakeInt(exitCode),
		},
	), nil
}

func starlarkListToStringSlice(list *starlark.List) []string {
	if list == nil {
		return []string{}
	}
	elems := make([]string, list.Len())
	for i := 0; i < list.Len(); i++ {
		elems[i] = list.Index(i).String()
	}
	return elems
}
