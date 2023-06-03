package starlarkos

import (
	"os"

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
	},
}

func getEnv(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name string
	if err := starlark.UnpackArgs("os.getenv", args, kwargs,
		"name", &name,
	); err != nil {
		return nil, err
	}
	if val, ok := os.LookupEnv(name); ok {
		return starlark.String(val), nil
	} else {
		return starlark.None, nil
	}
}
