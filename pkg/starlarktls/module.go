package starlarktls

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module starlarktls is a Starlark module of os-related functions and constants.
// The module defines the following functions:
//
//	getenv(name) - Gets the environment variable having the given name as a string, or None if not exists.
var Module = &starlarkstruct.Module{
	Name: "tls",
	Members: starlark.StringDict{
		"Config":         starlark.NewBuiltin("tls.Config", newConfig),
		"Certificate":    starlark.NewBuiltin("tls.Certificate", newCertificate),
		"ClientAuthType": clientAuthType,
	},
}
