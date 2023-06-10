package starlarkprocess

import (
	"os"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module process is a Starlark module of process-related functions and constants. The
// module defines the following functions:
func NewModule() *starlarkstruct.Module {
	executable, _ := os.Executable()
	return &starlarkstruct.Module{
		Name: "process",
		Members: starlark.StringDict{
			"executable": starlark.String(executable),
			"run":        starlark.NewBuiltin("process.run", run),
		},
	}
}
