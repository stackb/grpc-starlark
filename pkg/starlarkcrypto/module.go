package starlarkcrypto

import (
	"github.com/stackb/grpc-starlark/pkg/starlarkutil"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module starlarkcrypto is a Starlark module of crypto-related functions and constants.
// The module defines the following functions:
var Module = &starlarkstruct.Module{
	Name: "crypto",
	Members: starlark.StringDict{
		"x509": starlarkstruct.FromStringDict(
			starlarkutil.Symbol("x509"),
			starlark.StringDict{
				"SystemCertPool": starlark.NewBuiltin("x509.SystemCertPool", newSystemCertPool),
				"CertPool":       starlark.NewBuiltin("x509.CertPool", newCertPool),
			},
		),
		"tls": starlarkstruct.FromStringDict(
			starlarkutil.Symbol("tls"),
			starlark.StringDict{
				"Config":         starlark.NewBuiltin("tls.Config", newConfig),
				"Certificate":    starlark.NewBuiltin("tls.Certificate", newCertificate),
				"ClientAuthType": clientAuthType,
			},
		),
	},
}
