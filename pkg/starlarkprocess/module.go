package starlarkprocess

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module process is a Starlark module of process-related functions and constants. The
// module defines the following functions:
//
//	sleep(millis) - Sleeps the current process for the given number of milliseconds
var Module = &starlarkstruct.Module{
	Name: "process",
	Members: starlark.StringDict{
		"run": starlark.NewBuiltin("process.run", run),
	},
}
