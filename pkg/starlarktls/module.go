package starlarktls

import (
	"github.com/stackb/grpc-starlark/pkg/starlarkutil"
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
		"x509": starlarkstruct.FromStringDict(
			starlarkutil.Symbol("x509"),
			starlark.StringDict{
				"SystemCertPool": starlark.NewBuiltin("x509.SystemCertPool", newSystemCertPool),
				"CertPool":       starlark.NewBuiltin("x509.CertPool", newCertPool),
			},
		),
		"Config":         starlark.NewBuiltin("tls.Config", newConfig),
		"Certificate":    starlark.NewBuiltin("tls.Certificate", newCertificate),
		"ClientAuthType": clientAuthType,
	},
}
