package starlarkcrypto

import (
	"crypto/tls"

	"github.com/stackb/grpc-starlark/pkg/starlarkutil"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var clientAuthType = starlarkstruct.FromStringDict(
	starlarkutil.Symbol("tls.ClientAuthType"),
	starlark.StringDict{
		"NONE":               starlark.MakeInt(int(tls.NoClientCert)),
		"REQUEST":            starlark.MakeInt(int(tls.RequestClientCert)),
		"REQUIRE_ANY":        starlark.MakeInt(int(tls.RequireAnyClientCert)),
		"REQUIRE_AND_VERIFY": starlark.MakeInt(int(tls.RequireAndVerifyClientCert)),
	},
)
