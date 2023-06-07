package starlarkos

import (
	"os"

	"go.starlark.net/starlark"
)

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
