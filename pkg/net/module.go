package net

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module net is a Starlark module of network-related functions and constants. The
// module defines the following functions:
//
//	Listener(name) - creates a new net.Listener
var Module = &starlarkstruct.Module{
	Name: "net",
	Members: starlark.StringDict{
		"Listener": starlark.NewBuiltin("Listener", newListener),
	},
}
