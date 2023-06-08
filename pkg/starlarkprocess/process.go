package starlarkprocess

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/stackb/grpc-starlark/pkg/starlarkutil"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func run(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var command string
	var argv starlark.List
	var env starlark.Dict

	if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
		"command", &command,
		"args?", &argv,
		"env?", &env,
	); err != nil {
		return nil, err
	}

	cmd := exec.Command(command, starlarkListToStringSlice(&argv)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// cmd.Dir = "."

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
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	return starlarkstruct.FromStringDict(
		starlarkutil.Symbol(fn.Name()),
		starlark.StringDict{
			"exit_code": starlark.MakeInt(exitCode),
			"error":     starlark.String(errMsg),
			"stdout":    starlark.String(stdout),
			"stderr":    starlark.String(stderr),
		},
	), nil
}

func starlarkListToStringSlice(list *starlark.List) []string {
	elems := make([]string, list.Len())
	for i := 0; i < list.Len(); i++ {
		elems[i] = list.Index(i).String()
	}
	return elems
}
