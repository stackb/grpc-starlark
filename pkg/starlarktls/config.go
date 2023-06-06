package starlarktls

import (
	"crypto/tls"
	"fmt"

	"github.com/stackb/grpc-starlark/pkg/starlarkutil"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type Config struct {
	*tls.Config
	*starlarkstruct.Struct
}

func newConfig(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var certificates *starlark.List
	var clientAuth tls.ClientAuthType

	if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
		"certificates?", &certificates,
		"client_auth?", &clientAuth,
	); err != nil {
		return nil, err
	}

	certs := make([]tls.Certificate, certificates.Len())
	for i := 0; i < certificates.Len(); i++ {
		value := certificates.Index(i)
		if cert, ok := value.(certificate); ok {
			certs = append(certs, *cert.Certificate)
		} else {
			return nil, fmt.Errorf("certificate list entry %d must be a tls.Certificate (got %T)", i, value)
		}
	}

	return Config{
		Config: &tls.Config{
			Certificates: certs,
			ClientAuth:   clientAuth,
		},
		Struct: starlarkstruct.FromStringDict(
			starlarkutil.Symbol("tls.Config"),
			starlark.StringDict{
				"certificates": certificates,
			},
		),
	}, nil
}
