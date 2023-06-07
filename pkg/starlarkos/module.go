package starlarkos

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module starlarkos is a Starlark module of os-related functions and constants.
// The module defines the following functions:
//
//	getenv(name) - Gets the environment variable having the given name as a string, or None if not exists.
var Module = &starlarkstruct.Module{
	Name: "os",
	Members: starlark.StringDict{
		"getenv": starlark.NewBuiltin("os.getenv", getEnv),
		"stdin":  starlark.NewBuiltin("os.stdin", getStdin),
	},
}
